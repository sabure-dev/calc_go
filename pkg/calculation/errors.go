package calculation

type CalculationError struct {
	Message string
	Details string
}

func (e *CalculationError) Error() string {
	if e.Details != "" {
		return e.Message + ": " + e.Details
	}
	return e.Message
}

func NewExpressionTooShortError() error {
	return &CalculationError{
		Message: "выражение слишком короткое: минимальная длина 3 символа",
	}
}

func NewInvalidOperatorPositionError() error {
	return &CalculationError{
		Message: "некорректное расположение операторов: выражение не может начинаться или заканчиваться знаком операции",
	}
}

func NewDivisionByZeroError() error {
	return &CalculationError{
		Message: "деление на ноль",
	}
}

func NewInvalidCharError(char rune) error {
	return &CalculationError{
		Message: "недопустимый символ в выражении",
		Details: string(char),
	}
}

func NewBracketsExpressionError(err error) error {
	return &CalculationError{
		Message: "ошибка в выражении внутри скобок",
		Details: err.Error(),
	}
}

func NewSubExpressionError(err error) error {
	return &CalculationError{
		Message: "ошибка в вычислении подвыражения",
		Details: err.Error(),
	}
}

func NewConsecutiveOperatorsError() error {
	return &CalculationError{
		Message: "некорректное выражение: последовательные операторы недопустимы",
	}
}
