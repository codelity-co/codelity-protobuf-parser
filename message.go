package parser

type Message struct {
	name string
	fields map[string]interface{}
}

func (m *Message) GetMessageName() string {
	return m.name
}

func (m *Message) GetFields() map[string]interface{} {
	return m.fields
}