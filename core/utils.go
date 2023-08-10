package core

// Apply applique la fonction `transformer` à tous les éléments d'un slice et retourne le tableau convertit
func Apply[I interface{}, O interface{}](values []I, transformer func(I) O) []O {
	var output = []O{}
	for _, current := range values {
		output = append(output, transformer(current))
	}
	return output
}
