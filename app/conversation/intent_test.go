package conversation

import (
	"reflect"
	"testing"
)

func TestIntent(t *testing.T) {
	intent := NewIntent("TEST", "Que pense[sz] TU DE SUBJECT")

	testcases := []struct {
		In       string
		Expected map[string]string
	}{
		{"Que pensez vous de la vie ?", map[string]string{"subject": "vie", "pronoun": "vous"}},
		{"Que penses tu de l'amour ?", map[string]string{"subject": "amour", "pronoun": "tu"}},
		{"Que vois tu de l'amour ?", nil},
	}

	for i, testcase := range testcases {
		got := intent.Process(testcase.In)
		if !reflect.DeepEqual(testcase.Expected, got) {
			t.Fatalf("testcase #%v failed: expected %v, got %v\n", i, testcase.Expected, got)
		}
	}
}

func TestIntentCollection(t *testing.T) {
	collectionDef := [][2]string{
		{"TEST1", "TU? 1G(Parler) DE SUBJECT"},
		{"TEST2", "Que 1G(penser) TU DE SUBJECT"},
		{"TEST3", "Qu'est ce que LE? SUBJECT"},
		{"TEST4", "ce que TU 1G(penser) DE SUBJECT"},
		{"TEST5", "Et DE SUBJECT"},
		{"Free2", "(?:1G(composer)|1G(raconter)) un poème"},
		{"Greeting1", "^(?P<greeting>Bonjour|Salut|Hi|Coucou|Hello)"},
	}

	testcases := []struct {
		In             string
		ExpectedIntent string
		Expected       map[string]string
	}{
		{"Que pensez vous de la vie ?", "TEST2", map[string]string{"subject": "vie", "pronoun": "vous"}},
		{"Parle moi d'amour", "TEST1", map[string]string{"subject": "amour"}},
		{"dis moi doux robot, parles moi d'entropie", "TEST1", map[string]string{"subject": "entropie"}},
		{"Qu'est ce que l'amour ?", "TEST3", map[string]string{"subject": "amour"}},
		{"Qu'est ce que Dieu ?", "TEST3", map[string]string{"subject": "Dieu"}},
		{"Dis moi ce que tu penses de Dieu ?", "TEST4", map[string]string{"subject": "Dieu", "pronoun": "tu"}},
		{"J'aimerais parler de Dieu", "TEST1", map[string]string{"subject": "Dieu"}},
		{"J'aimerais que tu me parles de Dieu", "TEST1", map[string]string{"subject": "Dieu", "pronoun": "tu"}},
		{"J'aimerais que vous me parliez de Dieu", "TEST1", map[string]string{"subject": "Dieu", "pronoun": "vous"}},
		{"Hello robot !", "Greeting1", map[string]string{"greeting": "Hello"}},
		{"Hello robot, racontes moi un poème", "Free2", map[string]string{}},
		{"si vous pouviez me raconter un poème", "Free2", map[string]string{}},
		{"raconter moi un poeme", "Free2", map[string]string{}},
		{"Et des animaux", "TEST5", map[string]string{"subject": "animaux"}},
	}

	collection := NewIntentCollection(collectionDef)

	for i, testcase := range testcases {
		gotIntent, gotResult := FindIntent(collection, testcase.In)

		if gotIntent != testcase.ExpectedIntent {
			t.Fatalf("testcase #%v failed: expected intent %v, got %v\n", i, testcase.ExpectedIntent, gotIntent)
		}

		if !reflect.DeepEqual(testcase.Expected, gotResult) {
			t.Fatalf("testcase #%v failed: expected %v, got %v\n", i, testcase.Expected, gotResult)
		}
	}
}
