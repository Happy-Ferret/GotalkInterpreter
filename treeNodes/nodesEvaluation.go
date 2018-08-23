package treeNodes

import (
	"github.com/SealNTibbers/GotalkInterpreter/scanner"
	"fmt"
	"strconv"
)

type Scope struct {
	variables  map[string]SmalltalkObjectInterface
	OuterScope *Scope
}

func (s *Scope) Initialize() *Scope {
	s.variables = make(map[string]SmalltalkObjectInterface)
	return s
}

func (s *Scope) SetVar(name string, value SmalltalkObjectInterface) SmalltalkObjectInterface {
	s.variables[name] = value
	return value
}

func (s *Scope) SetStringVar(name string, value string) SmalltalkObjectInterface {
	smValue := NewSmalltalkString(value)
	s.variables[name] = smValue
	return smValue
}

func (s *Scope) SetNumberVar(name string, value float64) SmalltalkObjectInterface {
	smValue := NewSmalltalkNumber(value)
	s.variables[name] = smValue
	return smValue
}

func (s *Scope) SetBoolVar(name string, value bool) SmalltalkObjectInterface {
	smValue := NewSmalltalkBoolean(value)
	s.variables[name] = smValue
	return smValue
}

func (s *Scope) FindValueByName(name string) (SmalltalkObjectInterface, bool) {
	value, ok := s.variables[name]
	return value, ok
}

func (s *Scope) GetVarValue(name string) SmalltalkObjectInterface {
	value, ok := s.variables[name]
	if ok {
		return value
	} else {
		if s.OuterScope != nil {
			return s.OuterScope.GetVarValue(name)
		} else {
			panic(`we do not have variable with "` + name + `" in this scope`)
		}
	}
}

func (message *MessageNode) Eval(scope *Scope) SmalltalkObjectInterface {
	receiver := message.receiver.Eval(scope)
	var argObjects []SmalltalkObjectInterface
	for _, each := range message.arguments {
		argument := each.Eval(scope)
		if argument == nil {
			each.Eval(scope)
			panic("message argument should not be nil (void)")
		}
		argObjects = append(argObjects, argument)
	}
	result, err := receiver.Perform(message.GetSelector(), argObjects)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (block *BlockNode) Eval(scope *Scope) SmalltalkObjectInterface {
	return &SmalltalkBlock{&SmalltalkObject{}, block, scope}
}

func (sequence *SequenceNode) Eval(scope *Scope) SmalltalkObjectInterface {
	var result SmalltalkObjectInterface
	for _, each := range sequence.statements {
		result = each.Eval(scope)
	}
	return result
}

func (assignment *AssignmentNode) Eval(scope *Scope) SmalltalkObjectInterface {
	// create entry in our scope with assignment.variable and assignment.value
	scope.SetVar(assignment.variable.GetName(), assignment.value.Eval(scope))
	// return value for assignment variable
	return assignment.variable.Eval(scope)
}

func (variable *VariableNode) Eval(scope *Scope) SmalltalkObjectInterface {
	// return value for variable
	smalltalkValue := scope.GetVarValue(variable.GetName())
	if smalltalkValue.TypeOf() == DEFERRED {
		return smalltalkValue.Value()
	} else {
		return smalltalkValue
	}
}

func (array *LiteralArrayNode) Eval(scope *Scope) SmalltalkObjectInterface {
	arr := new(SmalltalkArray)
	for _, each := range array.contents {
		value := each.Eval(scope)
		arr.array = append(arr.array, value)
	}
	return arr
}

func (literalValue *LiteralValueNode) Eval(scope *Scope) SmalltalkObjectInterface {
	switch typeOfLiteral := literalValue.GetTypeOfToken(); typeOfLiteral {
	case scanner.NUMBER:
		{
			number, err := strconv.ParseFloat(literalValue.GetValue(), 64)
			if err == nil {
				object := new(SmalltalkNumber)
				object.SetValue(number)
				return object
			} else {
				return nil
			}
		}
	case scanner.STRING:
		{
			object := new(SmalltalkString)
			object.SetValue(literalValue.GetValue())
			return object
		}
	case scanner.BOOLEAN:
		{
			object := new(SmalltalkBoolean)
			object.SetValue(literalValue.GetValue() == "true")
			return object
		}
	default:
		return nil
	}
	return nil
}
