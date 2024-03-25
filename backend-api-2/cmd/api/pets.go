package main

type Pet struct {
	Nome   string `json:"nome"`
	Raca   string `json:"raca"`
	Tipo   string `json:"tipo"`
	Idade  string `json:"idade"`
	Genero string `json:"genero"`
}

type PetList struct {
	PetsSlice []Pet `json:"pets"`
}

var pets = []Pet{
	Pet{
		Nome:   "Tico",
		Raca:   "pinscher",
		Tipo:   "cachorro",
		Idade:  "8",
		Genero: "macho",
	},
	Pet{
		Nome:   "Teco",
		Raca:   "pinscher",
		Tipo:   "cachorro",
		Idade:  "8",
		Genero: "macho",
	},
	Pet{
		Nome:   "Nami",
		Raca:   "desconhecido",
		Tipo:   "gato",
		Idade:  "4",
		Genero: "femea",
	},
}
