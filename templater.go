// obsługa szablonów, ponieważ większość wykorzystuje w większości projektów
// więc jest to próba stworzenia biblioteki funkcji szablonów ogólnego zastosowania
package templater

import "fmt"

// Zarejestruj tylko gdy nie istnieje lub zwróć błąd
func Register(name string, f interface{}) error {
	if !Exist(name) {
		Overwrite(name, f)
		return nil
	}

	return fmt.Errorf("templater: Funkcja jest już zarejestrowana")
}

func Exist(name string) bool {
	return FindByName(name) != nil
}
func FindByName(name string) interface{} {
	if val, found := functions[name]; found {
		return val
	}
	return nil
}

func Overwrite(name string, f interface{}) {
	functions[name] = f
}

func DeleteByName(name string) {
	if Exist(name) {
		delete(functions, name)
	}
}

var functions = map[string]interface{}{}

