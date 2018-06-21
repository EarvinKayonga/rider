package domain

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/EarvinKayonga/rider/storage"
)

// PopulateDatabase put stuff in Database if empty
func PopulateDatabase(ctx context.Context, db storage.BikeStore) error {
	bikes, err := db.ListAllBikes(ctx, 1)
	if err != nil {
		return errors.Wrap(err, "an error occured while listing bikes")
	}

	if len(bikes) == 1 {
		return nil
	}

	toInsert, err := fromRegistry()
	if err != nil {
		return errors.Wrap(err, "an error occured while reading bikesRegistry")
	}

	_, err = db.CreateBikes(ctx, toInsert)
	if err != nil {
		return errors.Wrap(err, "an error occured while writing bikes to database")
	}

	return nil
}

func fromRegistry() ([]storage.Bike, error) {
	array := []jsonBike{}
	err := json.Unmarshal([]byte(bikesRegistry), &array)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while unmarshalling bikeRegistry")
	}

	bikes := []storage.Bike{}
	for index, b := range array {
		bike := storage.Bike{
			ID:        int64(index),
			PublicID:  b.ID,
			Status:    1,
			Latitude:  b.Location.Coordinates[0],
			Longitude: b.Location.Coordinates[1],
		}

		bikes = append(bikes, bike)
	}

	return bikes, nil
}

type jsonBike struct {
	ID       string `json:"id"`
	Location struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"location"`
}

const bikesRegistry = `
[{
  "id": "bb3398hl52n3nnikkt90",
  "location": {
    "type": "Point",
    "coordinates": [
      2.299104,
      48.854996
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkt9g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.34467,
      48.88958
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkta0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.293129,
      48.837738
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktag",
  "location": {
    "type": "Point",
    "coordinates": [
      2.278614,
      48.86362
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktb0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.396762,
      48.866942
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktbg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.399751,
      48.847796
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktc0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.385684,
      48.845773
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktcg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.337503,
      48.85047
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktd0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.392834,
      48.833481
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktdg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.368683,
      48.839103
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkte0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.369096,
      48.821439
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkteg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.388996,
      48.879582
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktf0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.349882,
      48.844236
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktfg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.364257,
      48.863829
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktg0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.279297,
      48.859422
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktgg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.325137,
      48.88447
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkth0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.28344,
      48.847355
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkthg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.369173,
      48.829103
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkti0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.388075,
      48.886486
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktig",
  "location": {
    "type": "Point",
    "coordinates": [
      2.363926,
      48.833381
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktj0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.319841,
      48.859007
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktjg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.332506,
      48.878644
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktk0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.326761,
      48.87523
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktkg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.33821,
      48.858982
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktl0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.316747,
      48.835304
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktlg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.300726,
      48.850856
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktm0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.335718,
      48.887093
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktmg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.334364,
      48.873815
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktn0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.366436,
      48.872052
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktng",
  "location": {
    "type": "Point",
    "coordinates": [
      2.355951,
      48.827642
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkto0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.355055,
      48.875175
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktog",
  "location": {
    "type": "Point",
    "coordinates": [
      2.365742,
      48.854261
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktp0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.382144,
      48.894707
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktpg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.292177,
      48.827215
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktq0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.377294,
      48.829708
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktqg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.305248,
      48.833165
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktr0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.284175,
      48.868481
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktrg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.344641,
      48.87747
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkts0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.289139,
      48.875661
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktsg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.401474,
      48.875494
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktt0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.292424,
      48.833276
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkttg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.306934,
      48.865746
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktu0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.309299,
      48.851149
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktug",
  "location": {
    "type": "Point",
    "coordinates": [
      2.385521,
      48.870823
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktv0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.375093,
      48.868503
    ]
  }
},
{
  "id": "bb3398hl52n3nnikktvg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.336137,
      48.841142
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku00",
  "location": {
    "type": "Point",
    "coordinates": [
      2.330605,
      48.833431
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku0g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.314107,
      48.871271
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku10",
  "location": {
    "type": "Point",
    "coordinates": [
      2.320767,
      48.834041
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku1g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.385836,
      48.861355
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku20",
  "location": {
    "type": "Point",
    "coordinates": [
      2.352543,
      48.836667
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku2g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.330273,
      48.853453
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku30",
  "location": {
    "type": "Point",
    "coordinates": [
      2.264677,
      48.848484
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku3g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.377014,
      48.851037
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku40",
  "location": {
    "type": "Point",
    "coordinates": [
      2.298551,
      48.833276
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku4g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.337994,
      48.897344
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku50",
  "location": {
    "type": "Point",
    "coordinates": [
      2.314617,
      48.853298
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku5g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.338176,
      48.876547
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku60",
  "location": {
    "type": "Point",
    "coordinates": [
      2.359788,
      48.841034
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku6g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.340001,
      48.829564
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku70",
  "location": {
    "type": "Point",
    "coordinates": [
      2.322392,
      48.867885
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku7g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.30093,
      48.882375
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku80",
  "location": {
    "type": "Point",
    "coordinates": [
      2.284103,
      48.839654
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku8g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.304319,
      48.841855
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku90",
  "location": {
    "type": "Point",
    "coordinates": [
      2.291223,
      48.866262
    ]
  }
},
{
  "id": "bb3398hl52n3nnikku9g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.288509,
      48.848221
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkua0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.299324,
      48.877074
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuag",
  "location": {
    "type": "Point",
    "coordinates": [
      2.347722,
      48.860718
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkub0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.358203,
      48.857515
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkubg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.319572,
      48.891979
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuc0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.285824,
      48.843514
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkucg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.293572,
      48.851452
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkud0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.298893,
      48.869343
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkudg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.324758,
      48.864033
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkue0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.320468,
      48.823953
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkueg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.331783,
      48.860418
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuf0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.294964,
      48.883996
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkufg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.307449,
      48.872943
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkug0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.322174,
      48.847282
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkugg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.314284,
      48.861677
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuh0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.402938,
      48.858617
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuhg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.352903,
      48.8962
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkui0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.305514,
      48.827331
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuig",
  "location": {
    "type": "Point",
    "coordinates": [
      2.350185,
      48.822701
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuj0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.310735,
      48.887797
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkujg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.339979,
      48.85336
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuk0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.339979,
      48.867195
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkukg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.316185,
      48.84902
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkul0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.27899,
      48.835448
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkulg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.273618,
      48.852818
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkum0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.28871,
      48.859553
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkumg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.364514,
      48.864287
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkun0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.348567,
      48.867846
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkung",
  "location": {
    "type": "Point",
    "coordinates": [
      2.308737,
      48.847166
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuo0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.291113,
      48.869698
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuog",
  "location": {
    "type": "Point",
    "coordinates": [
      2.395731,
      48.839195
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkup0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.321869,
      48.842392
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkupg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.328059,
      48.892803
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuq0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.296195,
      48.843669
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuqg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.364355,
      48.878944
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkur0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.332149,
      48.890144
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkurg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.374876,
      48.86913
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkus0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.355475,
      48.887239
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkusg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.304181,
      48.889763
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkut0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.340975,
      48.835603
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkutg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.372847,
      48.89136
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuu0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.304595,
      48.862975
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuug",
  "location": {
    "type": "Point",
    "coordinates": [
      2.33934,
      48.847405
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuv0",
  "location": {
    "type": "Point",
    "coordinates": [
      2.333369,
      48.847773
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkuvg",
  "location": {
    "type": "Point",
    "coordinates": [
      2.3686,
      48.88411
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkv00",
  "location": {
    "type": "Point",
    "coordinates": [
      2.275757,
      48.84293
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkv0g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.30631,
      48.877965
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkv10",
  "location": {
    "type": "Point",
    "coordinates": [
      2.315346,
      48.879972
    ]
  }
},
{
  "id": "bb3398hl52n3nnikkv1g",
  "location": {
    "type": "Point",
    "coordinates": [
      2.374626,
      48.858675
    ]
  }
}
]
`
