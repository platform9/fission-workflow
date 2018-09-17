package invocation

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/fission/fission-workflows/pkg/api"
	"github.com/fission/fission-workflows/pkg/controller"
	"github.com/fission/fission-workflows/pkg/controller/expr"
	"github.com/fission/fission-workflows/pkg/scheduler"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/fission/fission-workflows/pkg/types/typedvalues"
	"github.com/golang/protobuf/ptypes"
)

//
// Invocation-specific rules
//

type EvalContext interface {
	controller.EvalContext
	Workflow() *types.Workflow
	Invocation() *types.WorkflowInvocation
}

type WfiEvalContext struct {
	controller.EvalContext
	wf  *types.Workflow
	wfi *types.WorkflowInvocation
}

func NewEvalContext(state *controller.EvalState, wf *types.Workflow, wfi *types.WorkflowInvocation) WfiEvalContext {
	return WfiEvalContext{
		EvalContext: controller.NewEvalContext(state),
		wf:          wf,
		wfi:         wfi,
	}
}

func (ec WfiEvalContext) Workflow() *types.Workflow {
	return ec.wf
}

func (ec WfiEvalContext) Invocation() *types.WorkflowInvocation {
	return ec.wfi
}

type RuleWorkflowIsReady struct {
}

func (wr *RuleWorkflowIsReady) Eval(cec controller.EvalContext) controller.Action {
	ec := EnsureInvocationContext(cec)
	wf := ec.Workflow()
	// Check if workflow is in the right state to use.
	if !wf.Status.Ready() {
		log.WithField("wf.status", wf.Status.Status).Error("Workflow is not ready yet.")
		return &controller.ActionSkip{} // TODO backoff action
	}
	return nil
}

type RuleSchedule struct {
	Scheduler     *scheduler.WorkflowScheduler
	InvocationAPI *api.Invocation
	FunctionAPI   *api.Task
	StateStore    *expr.Store
}

func (sf *RuleSchedule) Eval(cec controller.EvalContext) controller.Action {
	ec := EnsureInvocationContext(cec)
	wf := ec.Workflow()
	wfi := ec.Invocation()
	// Request a execution plan from the Scheduler
	schedule, err := sf.Scheduler.Evaluate(&scheduler.ScheduleRequest{
		Invocation: wfi,
		Workflow:   wf,
	})
	if err != nil {
		return nil
	}

	// Execute the actions as specified in the execution plan
	var actions []controller.Action
	for _, a := range schedule.Actions {
		switch a.Type {
		case scheduler.ActionType_ABORT:
			invokeAction := &scheduler.AbortAction{}
			err := ptypes.UnmarshalAny(a.Payload, invokeAction)
			if err != nil {
				log.Errorf("Failed to unpack Scheduler action: %v", err)
			}
			return &ActionFail{
				API:          sf.InvocationAPI,
				InvocationID: wfi.ID(),
				Err:          errors.New(invokeAction.Reason),
			}
		case scheduler.ActionType_INVOKE_TASK:
			invokeAction := &scheduler.InvokeTaskAction{}
			err := ptypes.UnmarshalAny(a.Payload, invokeAction)
			if err != nil {
				log.Errorf("Failed to unpack Scheduler action: %v", err)
			}
			actions = append(actions, &ActionInvokeTask{
				Wf:         wf,
				ec:         ec.EvalState(),
				Wfi:        wfi,
				API:        sf.FunctionAPI,
				Task:       invokeAction,
				StateStore: sf.StateStore,
			})
		default:
			log.Warnf("Unknown Scheduler action: '%v'", a)
		}
	}
	return &controller.MultiAction{Actions: actions}
}

type RuleCheckIfCompleted struct {
	InvocationAPI *api.Invocation
}

func (cc *RuleCheckIfCompleted) Eval(cec controller.EvalContext) controller.Action {
	var (
		err   error
		ec    = EnsureInvocationContext(cec)
		wf    = ec.Workflow()
		wfi   = ec.Invocation()
		tasks = types.GetTasks(wf, wfi)
	)

	log.WithFields(logrus.Fields{
		"wfi-id": wfi.ID(),
		"wf-id":  wf.ID(),
	}).Debugf("check if all tasks are finished.")

	// Check if the workflow invocation is complete
	for id := range tasks {
		t, ok := wfi.Status.Tasks[id]
		if !ok || !t.Status.Finished() {
			return nil
		}
		if t.Status.Status != types.TaskInvocationStatus_SUCCEEDED {
			err = cc.InvocationAPI.Fail(wfi.ID(), errors.New("not all tasks succeeded"))
			if err != nil {
				return &controller.ActionError{
					Err: err,
				}
			}
		}
	}

	var finalOutput *types.TypedValue
	if len(wf.Spec.OutputTask) != 0 {
		finalOutput = typedvalues.ResolveTaskOutput(wf.Spec.OutputTask, wfi)
	}

	log.WithFields(logrus.Fields{
		"wfi-id":       wfi.ID(),
		"wf-id":        wf.ID(),
		"final-output": finalOutput,
	}).Debugf("all tasks finished.")

	// TODO extract to action
	err = cc.InvocationAPI.Complete(wfi.ID(), finalOutput)
	if err != nil {
		return &controller.ActionError{
			Err: err,
		}
	}

	return nil
}

func EnsureInvocationContext(cec controller.EvalContext) EvalContext {
	ec, ok := cec.(EvalContext)
	if !ok {
		panic("invalid evaluation context")
	}
	return ec
}
