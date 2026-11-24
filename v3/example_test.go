package audiometa_test

import (
	"io"

	"github.com/gcottom/audiometa/v3"
)

func ExampleOpenMetadata() {
	update := func(input io.ReadSeeker, output io.Writer) error {
		tag, err := audiometa.OpenMetadata(input)
		if err != nil {
			return err
		}
		if err := tag.SetTitle("A new title"); err != nil {
			return err
		}
		if err := tag.SetField(audiometa.Artist,
			audiometa.TextValue("First artist"),
			audiometa.TextValue("Second artist"),
		); err != nil {
			return err
		}
		if err := tag.SetCustom("COLLECTION", "Archive", "Favorites"); err != nil {
			return err
		}
		return tag.Save(output)
	}
	_ = update
}
