# structiterator

## Validator

Validator is a library to provide a way to put validation logic for struct as tag. This tag can be seen as a metadata for a struct.

Validator is using tag valid to inject the validation logic. Like the following example:
```
type App struct {
	Id             uint   `valid:"funcVal:Required"`
	Name           string `valid:"funcVal:Required"`
	Status         string
	ApprovalReason string `valid:"funcVal:CondRequired,compareKey:Status,compareValue:approved|rejected"`
}
```

For actually executing the validation logic, put the following lines in your code:

1. First you need to instantiate the validator object.

```
validtr := NewValidStruct()
```

2. Then execute ```validtr.Valid``` to run the intended validation logic.

```
app := App{}
app.Id = 2
app.Name = "First Test"
app.Status = "rejected"
validtr.Valid(app)
```
