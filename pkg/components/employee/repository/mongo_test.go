package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/migotom/cell-centre-services/pkg/entities"
	"github.com/migotom/cell-centre-services/pkg/helpers"
	"github.com/migotom/cell-centre-services/pkg/pb"
)

var (
	db *mongo.Database
)

func TestMain(m *testing.M) {
	var (
		code  = 1
		err   error
		purge func() error
	)

	db, purge, err = helpers.DBmock()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		code = m.Run()
	}

	purge()
	os.Exit(code)
}

func TestGet(t *testing.T) {
	cases := []struct {
		Name         string
		Filter       pb.EmployeeFilter
		Expectations func(*entities.Employee, error)
	}{
		{
			Name: "correct ID",
			Filter: pb.EmployeeFilter{
				Id: "5d3783ee28ae9468bc528906",
			},
			Expectations: func(employee *entities.Employee, err error) {
				employeeID, _ := primitive.ObjectIDFromHex("5d3783ee28ae9468bc528906")
				roleID, _ := primitive.ObjectIDFromHex("5d377ff93c9e1413c8c29e4b")
				now, _ := time.Parse(time.RFC3339, "2019-07-11T19:46:44Z")

				assert.NoError(t, err)
				assert.NotNil(t, employee)
				assert.Equal(t, &entities.Employee{
					ID:        employeeID,
					Email:     "admin@page.com",
					Password:  "$2a$04$FLsS84tlhkObb/m.ECwYdOXMkcq6lp1Yjrr5tptQTSf8xN3i2Y5wa",
					Name:      "John Doe",
					CreatedAt: &now,
					UpdatedAt: &now,
					Roles: []entities.Role{
						{
							ID:   roleID,
							Name: "admin",
						},
					},
				}, employee)
			},
		},
		{
			Name: "invalid ID",
			Filter: pb.EmployeeFilter{
				Id: "000083ee28ae9468bc528906",
			},
			Expectations: func(employee *entities.Employee, err error) {
				assert.EqualError(t, err, "mongo: no documents in result")
				assert.Nil(t, employee)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			repo := NewEmployeeRepository(db)
			tc.Expectations(repo.Get(context.Background(), &tc.Filter))
		})
	}
}
