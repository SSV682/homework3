package domain

type Status string

const (
	Success          Status = "success"
	Canceled         Status = "canceled"
	Created          Status = "created"
	PaymentPending   Status = "payment_pending"
	PaymentRejecting Status = "payment_rejecting"
	PaymentRejected  Status = "payment_rejected"
	PaymentApproved  Status = "payment_approved"
	StockPending     Status = "stock_pending"
	StockApproved    Status = "stock_approved"
	StockRejected    Status = "stock_reject"
	Canceling        Status = "canceling"
)

type Action string

const (
	NextStep Action = "next_step"
	Inaction Action = "inaction"
	End      Action = "end"
	Retry    Action = "retry"
)

type state interface {
	NextState(command OrderCommand) Step
}

type Step struct {
	Command Command
	Status  Status
	Action  Action
}

type Saga struct {
	CreatedState              state
	PaymentPendingState       state
	StockPendingState         state
	SuccessState              state
	PaymentRejectPendingState state
	StockRejectPendingState   state

	currentState state
	order        *Order
}

func (s *Saga) NextState(command OrderCommand) Step {
	step := s.currentState.NextState(command)

	switch step.Action {
	case NextStep:
		step.Command.Order = s.order
	default:
	}

	return step
}

func NewSaga(order *Order, paymentTopicName, stockTopicName string) *Saga {
	saga := Saga{
		order: order,
	}
	if order.Status == Created {
		saga.currentState = saga.CreatedState
	}

	if order.Status == Success {
		saga.currentState = saga.SuccessState
	}

	saga.CreatedState = NewCreatedState(&saga, paymentTopicName)
	saga.PaymentPendingState = NewPaymentPendingState(&saga, stockTopicName, paymentTopicName)
	saga.StockPendingState = NewStockPendingState(&saga, paymentTopicName, stockTopicName)
	saga.SuccessState = NewSuccessState(&saga, stockTopicName)
	saga.PaymentRejectPendingState = NewPaymentRejectPendingState(&saga, paymentTopicName)
	saga.StockRejectPendingState = NewStockRejectPendingState(&saga, paymentTopicName, stockTopicName)

	return &saga
}

type CreatedState struct {
	saga      *Saga
	topicName string
}

func (s *CreatedState) NextState(command OrderCommand) Step {
	if command.Status == Created {
		s.saga.currentState = s.saga.PaymentPendingState

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

func (s *PaymentPendingState) NextState(command OrderCommand) Step {
	switch command.Status {
	case PaymentApproved:
		s.saga.currentState = s.saga.StockPendingState
		return Step{
			Command: Command{
				Topic:       s.approveTopicName,
				CommandType: Approve,
			},
			Status: PaymentApproved,
			Action: NextStep,
		}
	case PaymentRejected:
		s.saga.currentState = nil
		return Step{
			Command: Command{},
			Status:  Canceled,
			Action:  End,
		}
	case Canceling:
		s.saga.currentState = s.saga.PaymentRejectPendingState
		return Step{
			Command: Command{
				Topic:       s.cancelingTopicName,
				CommandType: Reject,
			},
			Status: Canceling,
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
	rejectTopicName    string
	cancelingTopicName string
}

func (s *StockPendingState) NextState(command OrderCommand) Step {
	switch command.Status {
	case StockApproved:
		s.saga.currentState = s.saga.SuccessState
		return Step{
			Command: Command{},
			Status:  Success,
			Action:  End,
		}
	case StockRejected:
		s.saga.currentState = s.saga.PaymentRejectPendingState
		return Step{
			Command: Command{
				Topic:       s.rejectTopicName,
				CommandType: Reject,
			},
			Status: StockRejected,
			Action: NextStep,
		}
	case Canceling:
		s.saga.currentState = s.saga.PaymentRejectPendingState
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

func NewStockPendingState(saga *Saga, rejectTopicName, cancelingTopicName string) *StockPendingState {
	return &StockPendingState{
		saga:               saga,
		rejectTopicName:    rejectTopicName,
		cancelingTopicName: cancelingTopicName,
	}
}

type StockRejectPendingState struct {
	saga           *Saga
	nextTopicName  string
	retryTopicName string
}

func (s *StockRejectPendingState) NextState(command OrderCommand) Step {
	if command.Status == StockRejected {
		s.saga.currentState = s.saga.PaymentRejectPendingState
		return Step{
			Command: Command{
				Topic:       s.nextTopicName,
				CommandType: Reject,
			},
			Status: StockRejected,
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

func (s *PaymentRejectPendingState) NextState(command OrderCommand) Step {
	if command.Status == PaymentRejected {
		s.saga.currentState = nil
		return Step{
			Command: Command{},
			Status:  PaymentRejected,
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

func (s *SuccessState) NextState(command OrderCommand) Step {
	if command.Status == Canceling {
		s.saga.currentState = s.saga.StockRejectPendingState
		return Step{
			Command: Command{
				Topic:       s.cancelingTopicName,
				CommandType: Reject,
			},
			Status: StockRejected,
			Action: End,
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
