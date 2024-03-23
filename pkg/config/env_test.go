package config

type MyItem struct {
	Name string
}
type MyObj struct {
	Demo        string
	Stuff       string
	Items       []string
	Strongs     []*MyItem
	MyMap       map[string]string
	MyMapStrong map[string]*MyItem
}

// func Test_It(t *testing.T) {
// 	obj := &MyObj{
// 		Demo:        "demo",
// 		Stuff:       "stuff",
// 		Items:       []string{"one", "two"},
// 		Strongs:     []*MyItem{{Name: "one"}, {Name: "two"}},
// 		MyMap:       map[string]string{"one": "one", "two": "two"},
// 		MyMapStrong: map[string]*MyItem{"one": {Name: "one"}, "two": {Name: "two"}},
// 	}

// 	list, err := EnvMarshal("", obj)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(list) != 10 {
// 		t.Fatal("expected 10, got", len(list))
// 	}

// 	if list[0].Key != "Demo" {
// 		t.Fatal("expected Demo, got", list[0].Key)
// 	}
// 	if list[0].Value != "demo" {
// 		t.Fatal("expected demo, got", list[0].Value)
// 	}
// 	if list[1].Key != "Stuff" {
// 		t.Fatal("expected Stuff, got", list[1].Key)
// 	}
// 	if list[1].Value != "stuff" {
// 		t.Fatal("expected stuff, got", list[1].Value)
// 	}
// 	if list[2].Key != "Items_0" {
// 		t.Fatal("expected Items_0, got", list[2].Key)
// 	}
// 	if list[2].Value != "one" {
// 		t.Fatal("expected one, got", list[2].Value)
// 	}
// 	if list[3].Key != "Items_1" {
// 		t.Fatal("expected Items_1, got", list[3].Key)
// 	}
// 	if list[3].Value != "two" {
// 		t.Fatal("expected two, got", list[3].Value)
// 	}
// 	if list[4].Key != "Strongs_0_Name" {
// 		t.Fatal("expected Strongs_0_Name, got", list[4].Key)
// 	}
// 	if list[4].Value != "one" {
// 		t.Fatal("expected one, got", list[4].Value)
// 	}
// 	if list[5].Key != "Strongs_1_Name" {
// 		t.Fatal("expected Strongs_1_Name, got", list[5].Key)
// 	}
// 	if list[5].Value != "two" {
// 		t.Fatal("expected two, got", list[5].Value)
// 	}
// 	if list[6].Key != "MyMap_One" {
// 		t.Fatal("expected MyMap_One, got", list[6].Key)
// 	}
// 	if list[6].Value != "one" {
// 		t.Fatal("expected one, got", list[6].Value)
// 	}
// 	if list[7].Key != "MyMap_Two" {
// 		t.Fatal("expected MyMap_Two, got", list[7].Key)
// 	}
// 	if list[7].Value != "two" {
// 		t.Fatal("expected two, got", list[7].Value)
// 	}
// 	if list[8].Key != "MyMapStrong_One_Name" {
// 		t.Fatal("expected MyMapStrong_One_Name, got", list[8].Key)
// 	}
// 	if list[8].Value != "one" {
// 		t.Fatal("expected one, got", list[8].Value)
// 	}
// 	if list[9].Key != "MyMapStrong_Two_Name" {
// 		t.Fatal("expected MyMapStrong_Two_Name, got", list[9].Key)
// 	}
// 	if list[9].Value != "two" {
// 		t.Fatal("expected two, got", list[9].Value)
// 	}
// }

// func Test_UUID(t *testing.T) {
// 	id := uuid.New()

// 	list, err := EnvMarshal("", &id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if list[0].Key != "" {
// 		t.Fatal("expected empty, got", list[0].Key)
// 	}
// 	if list[0].Value != id.String() {
// 		t.Fatalf("expected %s, got %s", id.String(), list[0].Value)
// 	}
// }
