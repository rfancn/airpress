package assembler

import "github.com/rfancn/airpress/injection"

func init() {
	injection.Provide(
		NewBasePostAssembler,
		NewPostAssembler,
		NewSheetAssembler,
		NewBaseCommentAssembler,
		NewPostCommentAssembler,
		NewJournalCommentAssembler,
		NewSheetCommentAssembler,
	)
}
