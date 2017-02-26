package initializer

import (
	"fmt"
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

var (
	ObjectClass = InitializeObjectClass()
	ClassClass  = InitializeClassClass()
)

var BuiltinGlobalMethods = map[string]*object.BuiltInMethod{
	"puts": {
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return object.NULL
			}
		},
		Des:  "Print arguments",
		Name: "puts",
	},
}

var BuiltinClassMethods = map[string]*object.BuiltInMethod{
	"new": {
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				class := receiver.(*object.Class)
				instance := InitializeInstance(class)
				initMethod := class.LookupInstanceMethod("initialize")

				if initMethod != nil {
					instance.InitializeMethod = initMethod.(*object.Method)
				}

				return instance
			}
		},
		Des:  "Initialize class's instance",
		Name: "new",
	},
}

func InitializeMainObject() *object.BaseObject {
	obj := &object.BaseObject{Class: ObjectClass, InstanceVariables: object.NewEnvironment()}
	scope := &object.Scope{Self: obj, Env: object.NewEnvironment()}
	obj.Scope = scope
	return obj
}

func InitializeObjectClass() *object.Class {
	name := &ast.Constant{Value: "Object"}
	class := &object.Class{Name: name, Class: ClassClass, SuperClass: ClassClass, InstanceMethods: object.NewEnvironment(), ClassMethods: object.NewEnvironment()}

	return class
}

func InitializeClassClass() *object.Class {
	instanceMethods := object.NewEnvironment()
	classMethods := object.NewEnvironment()

	for key, value := range BuiltinGlobalMethods {
		instanceMethods.Set(key, value)
		classMethods.Set(key, value)
	}

	for key, value := range BuiltinClassMethods {
		classMethods.Set(key, value)
	}

	name := &ast.Constant{Value: "Class"}
	class := &object.Class{Name: name, InstanceMethods: instanceMethods, ClassMethods: classMethods}

	return class
}

func InitializeClass(name *ast.Constant, scope *object.Scope) *object.Class {
	class := &object.Class{Name: name, ClassMethods: object.NewEnvironment(), InstanceMethods: ClassClass.InstanceMethods, Class: ClassClass, SuperClass: ClassClass}
	classScope := &object.Scope{Self: class, Env: object.NewClosedEnvironment(scope.Env)}
	class.Scope = classScope

	return class
}

func InitializeInstance(c *object.Class) *object.BaseObject {
	instance := &object.BaseObject{Class: c, InstanceVariables: object.NewEnvironment()}

	return instance
}