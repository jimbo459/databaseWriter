package main

import (
	"database/sql"
	"errors"
)

type product struct {
	ID    int     `json:"ID"`
	Name  string  `json:"Name"`
	Price float64 `json:"Price"`
}

func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1", p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", p.Name, p.Price, p.ID)
	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *product) createProduct(db *sql.DB) error {
	return db.QueryRow("INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id", p.Name, p.Price).Scan(&p.ID)
}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	return nil, errors.New("Not implemented")
}
