package structures

func BottomUp[T any](
	start T, next func(current T) (proceeding T, exists bool), yield func(current T) bool,
) {
	var recursive func(element T)
	recursive = func(element T) {
		if yield(element) {
			if proceeding, exists := next(element); exists {
				recursive(proceeding)
			}
		}
	}
	recursive(start)
}
