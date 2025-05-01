package db

import (
	"fmt"
	"math/rand"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

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

var titles = []string{
	"Complete project documentation",
	"Review code changes",
	"Test API endpoints",
	"Update dependencies",
	"Fix UI bugs",
	"Write unit tests",
	"Optimize database queries",
	"Prepare presentation",
}

// func Seed(store store.Storage) {
// 	ctx := context.Background()

// 	users := generateUsers(10)
// 	for _, user := range users {
// 		if err := store.Users.Create(ctx, user); err != nil {
// 			log.Println("Error creating user:", err)
// 			return
// 		}
// 	}

// 	tasks := generateTasks(20, users)
// 	if err := store.Tasks.CreateTasks(ctx, tasks); err != nil {
// 		log.Println("Error creating tasks:", err)
// 		return
// 	}

// }

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
			//Password:    "123",
		}
	}
	return users
}

func generateTasks(num int, users []*store.User) []*store.Task {
	tasks := make([]*store.Task, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		tasks[i] = &store.Task{
			UserID:     user.ID,
			Title:      titles[rand.Intn(len(titles))],
			IsOptional: rand.Intn(2) == 1,
		}
	}
	return tasks
}
