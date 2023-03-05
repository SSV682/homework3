package domain

import "github.com/labstack/gommon/log"

type Status string

const (
	Success           Status = "success"
	Canceled          Status = "canceled"
	Created           Status = "created"
	PaymentPending    Status = "payment_pending"
	PaymentRejecting  Status = "payment_rejecting"
	PaymentRejected   Status = "payment_rejected"
	PaymentApproved   Status = "payment_approved"
	StockPending      Status = "stock_pending"
	StockApproved     Status = "stock_approved"
	StockRejected     Status = "stock_rejected"
	StockRejecting    Status = "stock_rejecting"
	DeliveryPending   Status = "delivery_pending"
	DeliveryApproved  Status = "delivery_approved"
	DeliveryRejected  Status = "delivery_rejected"
	DeliveryRejecting Status = "delivery_rejecting"
	Canceling         Status = "canceling"
)

type Action string

const (
	NextStep Action = "next_step"
	Inaction Action = "inaction"
	End      Action = "end"
	Retry    Action = "retry"
)

type state interface {
	NextStep(command OrderCommand) Step
}

type Step struct {
	Command Command
	Status  Status
	Action  Action
}

type Saga struct {
	createdState               state
	paymentPendingState        state
	stockPendingState          state
	deliveryPendingState       state
	successState               state
	paymentRejectPendingState  state
	stockRejectPendingState    state
	deliveryRejectPendingState state
	canceledState              state
	currentState               state
	order                      *Order
}

func (s *Saga) SetState(state state) {
	s.currentState = state
}

func (s *Saga) NextState(command OrderCommand) Step {
	step := s.currentState.NextStep(command)
	switch step.Action {
	case NextStep:
		step.Command.Order = s.order
	default:
	}

	return step
}

func NewSaga(order *Order, paymentTopicName, stockTopicName, deliveryTopicName string) *Saga {
	saga := Saga{
		order: order,
	}

	saga.createdState = NewCreatedState(&saga, paymentTopicName)
	saga.paymentPendingState = NewPaymentPendingState(&saga, stockTopicName, paymentTopicName)
	saga.stockPendingState = NewStockPendingState(&saga, deliveryTopicName, paymentTopicName, stockTopicName)
	saga.deliveryPendingState = NewDeliveryPendingState(&saga, stockTopicName, deliveryTopicName)
	saga.successState = NewSuccessState(&saga, deliveryTopicName)
	saga.paymentRejectPendingState = NewPaymentRejectPendingState(&saga, paymentTopicName)
	saga.stockRejectPendingState = NewStockRejectPendingState(&saga, paymentTopicName, stockTopicName)
	saga.deliveryRejectPendingState = NewDeliveryRejectPendingState(&saga, stockTopicName, deliveryTopicName)
	saga.canceledState = NewCanceledState(&saga)

	if order.Status == Created {
		saga.currentState = saga.createdState
	}

	if order.Status == Success {
		saga.currentState = saga.successState
	}

	log.Infof("saga for orderID %d created: %v", order.ID, saga)

	return &saga
}

type CreatedState struct {
	saga      *Saga
	topicName string
}

func (s *CreatedState) NextStep(command OrderCommand) Step {
	if command.Status == Created {
		s.saga.SetState(s.saga.paymentPendingState)
		return Step{
			Command: Command{
				Topic:       s.topicName,
				CommandType: Approve,
			},
			Status: PaymentPending,
			Action: NextStep,
		}
	}

	return Step{
		Command: Command{},
		Status:  Canceled,
		Action:  End,
	}
}

func NewCreatedState(saga *Saga, topicName string) *CreatedState {
	return &CreatedState{
		saga:      saga,
		topicName: topicName,
	}
}

type PaymentPendingState struct {
	saga               *Saga
	approveTopicName   string
	cancelingTopicName string
}

func (s *PaymentPendingState) NextStep(command OrderCommand) Step {
	switch command.Status {
	case PaymentApproved:
		s.saga.SetState(s.saga.stockPendingState)
		return Step{
			Command: Command{
				Topic:       s.approveTopicName,
				CommandType: Approve,
			},
			Status: StockPending,
			Action: NextStep,
		}
	case PaymentRejected:
		s.saga.SetState(s.saga.canceledState)
		return Step{
			Command: Command{},
			Status:  Canceled,
			Action:  End,
		}
	case Canceling:
		s.saga.SetState(s.saga.paymentRejectPendingState)
		return Step{
			Command: Command{
				Topic:       s.cancelingTopicName,
				CommandType: Reject,
			},
			Status: Canceled,
			Action: NextStep,
		}
	}

	return Step{
		Command: Command{},
		Status:  Canceled,
		Action:  Inaction,
	}
}

func NewPaymentPendingState(saga *Saga, approveTopicName, cancelingTopicName string) *PaymentPendingState {
	return &PaymentPendingState{
		saga:               saga,
		approveTopicName:   approveTopicName,
		cancelingTopicName: cancelingTopicName,
	}
}

type StockPendingState struct {
	saga               *Saga
	approvedTopicName  string
	rejectTopicName    string
	cancelingTopicName string
}

func (s *StockPendingState) NextStep(command OrderCommand) Step {
	switch command.Status {
	case StockApproved:
		s.saga.SetState(s.saga.deliveryPendingState)
		return Step{
			Command: Command{
				Topic:       s.approvedTopicName,
				CommandType: Approve,
			},
			Status: DeliveryPending,
			Action: NextStep,
		}
	case StockRejected:
		s.saga.SetState(s.saga.paymentRejectPendingState)
		return Step{
			Command: Command{
				Topic:       s.rejectTopicName,
				CommandType: Reject,
			},
			Status: StockRejected,
			Action: NextStep,
		}
	case Canceling:
		s.saga.SetState(s.saga.paymentRejectPendingState)
		return Step{
			Command: Command{
				Topic:       s.cancelingTopicName,
				CommandType: Reject,
			},
			Status: PaymentRejecting,
			Action: NextStep,
		}
	}

	return Step{
		Command: Command{},
		Status:  StockPending,
		Action:  Inaction,
	}
}

func NewStockPendingState(saga *Saga, approvedTopicName, rejectTopicName, cancelingTopicName string) *StockPendingState {
	return &StockPendingState{
		saga:               saga,
		approvedTopicName:  approvedTopicName,
		rejectTopicName:    rejectTopicName,
		cancelingTopicName: cancelingTopicName,
	}
}

type StockRejectPendingState struct {
	saga           *Saga
	nextTopicName  string
	retryTopicName string
}

func (s *StockRejectPendingState) NextStep(command OrderCommand) Step {
	if command.Status == StockRejected {
		s.saga.SetState(s.saga.paymentRejectPendingState)
		return Step{
			Command: Command{
				Topic:       s.nextTopicName,
				CommandType: Reject,
			},
			Status: PaymentRejecting,
			Action: NextStep,
		}
	}

	return Step{
		Command: Command{},
		Status:  StockRejected,
		Action:  Inaction,
	}
}

func NewStockRejectPendingState(saga *Saga, nextTopicName, retryTopicName string) *StockRejectPendingState {
	return &StockRejectPendingState{
		saga:           saga,
		nextTopicName:  nextTopicName,
		retryTopicName: retryTopicName,
	}
}

type PaymentRejectPendingState struct {
	saga           *Saga
	retryTopicName string
}

func (s *PaymentRejectPendingState) NextStep(command OrderCommand) Step {
	if command.Status == PaymentRejected {
		s.saga.SetState(s.saga.canceledState)
		return Step{
			Command: Command{},
			Status:  Canceled,
			Action:  End,
		}
	}

	return Step{
		Command: Command{},
		Status:  PaymentPending,
		Action:  Inaction,
	}
}

func NewPaymentRejectPendingState(saga *Saga, retryTopicName string) *PaymentRejectPendingState {
	return &PaymentRejectPendingState{
		saga:           saga,
		retryTopicName: retryTopicName,
	}
}

type SuccessState struct {
	saga               *Saga
	cancelingTopicName string
}

func (s *SuccessState) NextStep(command OrderCommand) Step {
	if command.Status == Canceling {
		s.saga.SetState(s.saga.deliveryRejectPendingState)
		return Step{
			Command: Command{
				Topic:       s.cancelingTopicName,
				CommandType: Reject,
			},
			Status: DeliveryRejecting,
			Action: NextStep,
		}
	}

	return Step{
		Command: Command{},
		Status:  Success,
		Action:  Inaction,
	}
}

func NewSuccessState(saga *Saga, cancelingTopicName string) *SuccessState {
	return &SuccessState{
		saga:               saga,
		cancelingTopicName: cancelingTopicName,
	}
}

type CanceledState struct {
	saga *Saga
}

func (s *CanceledState) NextStep(_ OrderCommand) Step {
	return Step{
		Command: Command{},
		Status:  Canceled,
		Action:  Inaction,
	}
}

func NewCanceledState(saga *Saga) *CanceledState {
	return &CanceledState{
		saga: saga,
	}
}

type DeliveryPendingState struct {
	saga               *Saga
	rejectTopicName    string
	cancelingTopicName string
}

func (s *DeliveryPendingState) NextStep(command OrderCommand) Step {
	switch command.Status {
	case DeliveryApproved:
		s.saga.SetState(s.saga.successState)
		return Step{
			Command: Command{},
			Status:  Success,
			Action:  End,
		}
	case DeliveryRejected:
		s.saga.SetState(s.saga.stockRejectPendingState)
		return Step{
			Command: Command{
				Topic:       s.rejectTopicName,
				CommandType: Reject,
			},
			Status: StockRejecting,
			Action: NextStep,
		}
	case Canceling:
		s.saga.SetState(s.saga.deliveryRejectPendingState)
		return Step{
			Command: Command{
				Topic:       s.cancelingTopicName,
				CommandType: Reject,
			},
			Status: DeliveryRejecting,
			Action: NextStep,
		}
	}

	return Step{
		Command: Command{},
		Status:  StockPending,
		Action:  Inaction,
	}
}

func NewDeliveryPendingState(saga *Saga, rejectTopicName, cancelingTopicName string) *DeliveryPendingState {
	return &DeliveryPendingState{
		saga:               saga,
		rejectTopicName:    rejectTopicName,
		cancelingTopicName: cancelingTopicName,
	}
}

type DeliveryRejectPendingState struct {
	saga           *Saga
	nextTopicName  string
	retryTopicName string
}

func (s *DeliveryRejectPendingState) NextStep(command OrderCommand) Step {
	if command.Status == DeliveryRejected {
		s.saga.SetState(s.saga.stockRejectPendingState)
		return Step{
			Command: Command{
				Topic:       s.nextTopicName,
				CommandType: Reject,
			},
			Status: StockRejecting,
			Action: NextStep,
		}
	}

	return Step{
		Command: Command{},
		Status:  DeliveryRejected,
		Action:  Inaction,
	}
}

func NewDeliveryRejectPendingState(saga *Saga, nextTopicName, retryTopicName string) *DeliveryRejectPendingState {
	return &DeliveryRejectPendingState{
		saga:           saga,
		nextTopicName:  nextTopicName,
		retryTopicName: retryTopicName,
	}
}
