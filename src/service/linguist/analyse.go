package linguist
import (
	"backend/src/service/linguist/data"
	tokenizer "backend/src/service/tokenizer"
	"bytes"
	"github.com/jbrukh/bayesian"
	"log"
	"math"
)

var classifier *bayesian.Classifier
var classifier_initialized bool = false

// Gets the baysian.Classifier which has been trained on programming language
// samples from github.com/github/linguist after running the generator
//
// See also cmd/generate-classifier
func getClassifier() *bayesian.Classifier {
	// NOTE(tso): this could probably go into an init() function instead
	// but this lazy loading approach works, and it's conceivable that the
	// analyse() function might not invoked in an actual runtime anyway
	if !classifier_initialized {
		data, err := data.Asset("classifier")
		if err != nil {
			log.Panicln(err)
		}
		reader := bytes.NewReader(data)
		classifier, err = bayesian.NewClassifierFromReader(reader)
		if err != nil {
			log.Panicln(err)
		}
		classifier_initialized = true
	}
	return classifier
}
func Analyse(contents []byte, hints []string) (language string) {
	document := tokenizer.Tokenize(contents)
	classifier := getClassifier()
	scores, idx, _ := classifier.LogScores(document)

	if len(hints) == 0 {
		return string(classifier.Classes[idx])
	}

	langs := map[string]struct{}{}
	for _, hint := range hints {
		langs[hint] = struct{}{}
	}

	best_score := math.Inf(-1)
	best_answer := ""

	for id, score := range scores {
		answer := string(classifier.Classes[id])
		if _, ok := langs[answer]; ok {
			if score >= best_score {
				best_score = score
				best_answer = answer
			}
		}
	}
	return best_answer
}
