package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

var password = "123123"

var usernames = []string{
	"john_doe",
	"jane_smith",
	"mike_j",
	"sarah_k",
	"alex_g",
	"emily_w",
	"david_m",
	"lisa_b",
	"ryan_t",
	"amy_c",
	"kevin_l",
	"olivia_n",
	"steve_r",
	"mia_j",
	"daniel_k",
}

var displayNames = []string{
	"John Doe",
	"Jane Smith",
	"Mike Johnson",
	"Sarah Williams",
	"Alex Garcia",
	"Emily Wilson",
	"David Martinez",
	"Lisa Brown",
	"Ryan Taylor",
	"Amy Clark",
	"Kevin Lee",
	"Olivia Nelson",
	"Steve Roberts",
	"Mia Jones",
	"Daniel Kim",
}

var taskTitles = []string{
	"Complete project documentation",
	"Review code changes",
	"Test API endpoints",
	"Update dependencies",
	"Fix UI bugs",
	"Write unit tests",
	"Optimize database queries",
	"Prepare presentation",
}

var goalTitles = []string{
	"Learn Go",
	"Complete 100 days of code challenge",
	"Build a personal project",
	"Master database design",
	"Improve system architecture skills",
	"Contribute to open source",
	"Learn cloud technologies",
	"Improve code review skills",
	"Attend 3 tech conferences this year",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(5)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	tasks := generateTasks(20, users)
	if err := store.Tasks.CreateTasks(ctx, tasks); err != nil {
		log.Println("Error creating tasks:", err)
		return
	}

	// goals := generateGoals(15, users)
	// if err := store.Goals.CreateGoals(ctx, goals); err != nil {
	// 	log.Println("Error creating goals:", err)
	// 	return
	// }

}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		var displayName *string
		if rand.Intn(5) > 0 { // 80% chance of having a name (4 out of 5)
			name := displayNames[i%len(displayNames)]
			displayName = &name
		}

		users[i] = &store.User{
			Username:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:       usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			DisplayName: displayName,
		}
		users[i].Password.Set(password)
	}
	return users
}

func generateTasks(num int, users []*store.User) []*store.Task {
	tasks := make([]*store.Task, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		tasks[i] = &store.Task{
			UserID:     user.ID,
			Title:      taskTitles[rand.Intn(len(taskTitles))],
			IsOptional: rand.Intn(2) == 1,
		}
	}
	return tasks
}

func generateGoals(num int, users []*store.User) []*store.Goal {
	goals := make([]*store.Goal, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		goals[i] = &store.Goal{
			UserID: user.ID,
			Title:  goalTitles[rand.Intn(len(goalTitles))],
		}
	}
	return goals
}
