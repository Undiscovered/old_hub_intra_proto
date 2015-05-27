package db

import (
    "intra-hub/models"
)

// Populate the database with default values.
func PopulateDatabase() {
    // Add Groups
    i, err := QueryGroup().PrepareInsert()
    if err != nil {
        panic(err)
    }
    for _, groupName := range models.EveryUserGroups {
        group := &models.Group{
            Name: groupName,
        }
        i.Insert(group)
    }
    if err := i.Close(); err != nil {
        panic(err)
    }
}
