package object

import (
	"fmt"
	"strconv"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	String() string
	IsTruthy() bool
}

const (
	INTEGER ObjectType = "INTEGER"
	FLOAT   ObjectType = "FLOAT"
	STRING  ObjectType = "STRING"
	BOOLEAN ObjectType = "BOOLEAN"
	NIL     ObjectType = "NIL"
	ERROR   ObjectType = "ERROR"
	RETURN  ObjectType = "RETURN"
)

type Integer struct {
	Value int
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

func (i *Integer) String() string {
	return strconv.Itoa(i.Value)
}

func (i *Integer) IsTruthy() bool {
	return i.Value != 0
}

func NewInteger(value int) *Integer {
	return &Integer{Value: value}
}

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType {
	return "FLOAT"
}

func (f *Float) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

func (f *Float) IsTruthy() bool {
	return f.Value != 0
}

func NewFloat(value float64) *Float {
	return &Float{Value: value}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (b *Boolean) String() string {
	return strconv.FormatBool(b.Value)
}

func (b *Boolean) IsTruthy() bool {
	return b.Value
}

type String struct {
	Value string
}

func NewString(value string) *String {
	return &String{Value: value}
}

func (s *String) Type() ObjectType {
	return STRING
}

func (s *String) String() string {
	return s.Value
}

func (s *String) IsTruthy() bool {
	return len(s.Value) != 0
}

func NewBoolean(value bool) *Boolean {
	return &Boolean{Value: value}
}

type Nil struct{}

func (n *Nil) Type() ObjectType {
	return NIL
}

func (n *Nil) String() string {
	return "nil"
}

func (n *Nil) IsTruthy() bool {
	return false
}

func NewNil() *Nil {
	return &Nil{}
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR
}

func (e *Error) String() string {
	return "ERROR: " + e.Message
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN
}

func (rv *ReturnValue) String() string {
	return rv.Value.String()
}

func (rv *ReturnValue) IsTruthy() bool {
	return rv.Value.IsTruthy()
}

func (e *Error) IsTruthy() bool {
	return false
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func IsNumber(obj Object) bool {
	return obj.Type() == INTEGER || obj.Type() == FLOAT
}

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Set(name string, value Object) {
	e.store[name] = value
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}
