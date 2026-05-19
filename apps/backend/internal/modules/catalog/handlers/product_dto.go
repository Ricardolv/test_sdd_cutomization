package handlers

type CreateProductRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Status          string  `json:"status"`
	AvailableOnline bool    `json:"availableOnline"`
	Featured        bool    `json:"featured"`
	AllowsPreOrder  bool    `json:"allowsPreOrder"`
}

type UpdateProductRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Status          string  `json:"status"`
	AvailableOnline bool    `json:"availableOnline"`
	Featured        bool    `json:"featured"`
	AllowsPreOrder  bool    `json:"allowsPreOrder"`
}

type ProductResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Status          string  `json:"status"`
	AvailableOnline bool    `json:"availableOnline"`
	Featured        bool    `json:"featured"`
	AllowsPreOrder  bool    `json:"allowsPreOrder"`
}

type ListResponse struct {
	Items []ProductResponse `json:"items"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}
