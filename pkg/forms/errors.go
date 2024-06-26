package forms

// Create a type for errors.
type errors map[string][]string;

// Add a method to the errors type that adds a new error to the map.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message);
}

func (e errors) Get(field string) string {
	es := e[field];
	if len(es) == 0 {
		return "";
	}
	return es[0];
}