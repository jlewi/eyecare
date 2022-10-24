package scraping

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func Test_HtmlToText(t *testing.T) {
	type testCase struct {
		fileName string
		expected string
	}

	cases := []testCase{
		{
			fileName: "about_me.html",
			expected: `ABOUT ME | website
ABOUT DR. KOUTOULAS
Dr. Kosta A. Koutoulas was born and raised in the Bay Area. After moving to Chicago to attend the Illinois College of Optometry, he witnessed firsthand the impact eye doctors can have on the lives of their patients. He was inspired by the mixture of mathematics and technology inherent to optometry work, as well as the opportunity to help people by restoring or improving their gift of vision.
In 2007, Dr. Koutoulas was thrilled to open his private optometry practice in Daly City. Today, he loves the time and flexibility that comes with running his own practice, giving his full attention to each patient, listening to their specific symptoms, and working with them to locate the source of an ocular problem.
Dr. Koutoulas believes that being psychologically and intellectually prepared is essential to the success of any treatment, so he works to educate patients about their condition, as well as the different treatment options that are available. Often times there is more than one solution to an ocular problem, and Dr. Koutoulas strives to work together with the patient in deciding the best course of action.
Using state-of-the-art technology such as Zeiss corneal topographers, which map the shape and condition of your cornea, Dr. Koutoulas tailors contact lenses to fit your eye’s unique shape and needs. The Humphrey Visual Field analyzer is a great way to diagnose eye conditions that might be affecting a patient’s optic nerve, such as glaucoma. These images, along with pictures taken with our in-office technology, can be shared via smartphone or email, allowing patients to see exactly what Dr. Koutoulas is seeing. More than just offering great treatment, he wants to help each patient become an active participant in their own care.
I want everyone to be confident and comfortable with the care they’re receiving. I hope patients feel educated when they leave their first appointment, like they’ve learned something new.
DR. KOSTA A. KOUTOULAS
PERSONAL INTERESTS
Dr. Koutoulas is an avid runner, participating regularly in half-marathons and Tough Mudders. He enjoys traveling, and he speaks three languages: Greek, English, and Spanish.
PERSONAL INTERESTS
Dr. Koutoulas is always looking for ways to expand his practice and offer more advanced treatments to his patients. He has been certified to treat patients suffering from glaucoma, a disease that impacts vision by damaging the optic nerve. He is also certified in corneal refractive therapy, a new, noninvasive method of vision correction that can improve sight with the use of glasses, contacts, or cornea surgery. He’s proud to be one of the only doctors in the country to offer this groundbreaking new method of treatment.
2001
UNDERGRADUATE
University of California, Los Angeles
BSc.Physiological Science
2006
2005-2006
clinical rotations
Daytona Beach Veterans Hospital
Ocular disease and primary eye care
San Diego Vision Development Center
Pediatric and adult vision therapy and neurological development
Lawndale Christian Health Center
Primary and urgent eye care
Illinois Eye Institute
Cornea and Contact Lens center
Center for Advanced Ophthalmic care
-Glaucoma section
Illinois College of Optometry
O.D.
Orthokeratology training and myopia reduction
MEMBERSHIPS
San Mateo Optometric Association
Daly City Host Lions club eye provider, Daly City
VA/QTC eye provider- Veteran's compensation
`,
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory; %v", err)
	}
	testDir := filepath.Join(cwd, "test_data")
	for _, c := range cases {
		t.Run(c.fileName, func(t *testing.T) {
			path := filepath.Join(testDir, c.fileName)
			b, err := os.Open(path)
			if err != nil {
				t.Fatalf("Failed to open file %v; error %v", path, err)
			}

			actual := TextFromHtml(b)
			fmt.Printf("%v\n", actual)
			if d := cmp.Diff(c.expected, actual); d != "" {
				t.Errorf("Unexpected diff:\n%v", d)
			}
		})
	}
}
