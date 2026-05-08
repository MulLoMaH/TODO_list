package domain

// Дженерик создает указатель на любой передаваемый тип данных
type Nullable[T any] struct {
	Value *T
	Set   bool
}

//Использование [передается нужный тип данных]
//Создает указатель на строку
// NullableString := Nullable[string]{}

//Создаем ограничение для дженерика что бы он мог принимать только строку
// type Speaker interface {
// 	Speak() string
// }

//Теперь дженерик принимает не все типы а только строку так как
// тип данных Speaker который хранит в себе строку
// type Nullable[T Speaker] struct {
// 	value *T
// 	Set bool
// }
