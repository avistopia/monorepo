package components

type Renderer[T any] interface {
	Render() T
}

func renderTable[T any](rows [][]Renderer[T]) [][]T {
	renderedTable := make([][]T, 0, len(rows))

	for _, row := range rows {
		renderedRow := make([]T, 0, len(row))

		for _, item := range row {
			renderedRow = append(renderedRow, item.Render())
		}

		renderedTable = append(renderedTable, renderedRow)
	}

	return renderedTable
}
