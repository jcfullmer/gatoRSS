package Config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	database "github.com/jcfullmer/gatoRSS/internal/database"
	"github.com/jcfullmer/gatoRSS/internal/rss"
)

type State struct {
	Db     *database.Queries
	Config *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func GetState() (*State, error) {
	config, err := Read()
	if err != nil {
		return &State{}, err
	}

	return &State{
		nil,
		&config,
	}, nil
}

func CreateCommand(args []string) Command {
	name := args[1]
	var argies []string
	if len(args) >= 2 {
		argies = args[2:]
	}
	return Command{
		Name: name,
		Args: argies,
	}
}

func (c *Commands) Run(s *State, cmd Command) error {
	if _, ok := c.Handlers[cmd.Name]; !ok {
		return fmt.Errorf("command %v not found\n", cmd.Name)
	}
	err := c.Handlers[cmd.Name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	if _, ok := c.Handlers[name]; ok {
		return fmt.Errorf("command already exists")
	}
	c.Handlers[name] = f
	return nil
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("a username is required")
		os.Exit(1)
	}
	u, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("User doesn't exist")
		os.Exit(1)
	}
	newName := cmd.Args[0]
	err = s.Config.SetUser(newName, u.ID)
	if err != nil {
		return err
	}
	fmt.Println("The user has been set.")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Please input a name to register.")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err == sql.ErrNoRows {
		dbUser := database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
		}

		user, err := s.Db.CreateUser(context.Background(), dbUser)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		Config.SetUser(*s.Config, dbUser.Name, dbUser.ID)
		fmt.Println("the user was created")
		fmt.Println(user)
		return nil
	} else if err != nil {
		fmt.Println("an error occured when getting user", err)
		os.Exit(1)
	} else {
		fmt.Println("user already exists")
		os.Exit(1)
	}
	return nil
}

func HandlerReset(s *State, _ Command) error {
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error when resetting databse: %v", err)
		os.Exit(1)
	}
	fmt.Println("Successfully reset database.")
	os.Exit(0)
	return nil
}

func HandlerUsers(s *State, _ Command) error {
	list, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error when Getting users: %v\n", err)
		os.Exit(1)
	}
	if len(list) == 0 {
		fmt.Println("No users in database.")
		os.Exit(1)
	}
	for _, user := range list {
		if user == s.Config.CurrentUserName {
			fmt.Printf("* %v (current)\n", user)
			continue
		}
		fmt.Printf("* %v\n", user)
	}
	os.Exit(0)
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Please input a time")
		os.Exit(1)
	}
	time_between_regs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		fmt.Println("error when parsing time between regs: ", err)
		os.Exit(1)
	}
	fmt.Println("Collecting feeds every ", time_between_regs)
	ticker := time.NewTicker(time_between_regs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		fmt.Println("not nough args, need name of feed and url")
		os.Exit(1)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	dbFeed, err := s.Db.CreateFeed(context.Background(), feed)

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    s.Config.CurrentUserID,
		FeedID:    dbFeed.ID,
	}
	feedRows, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		fmt.Println("Error when following feed", err)
		os.Exit(1)
	}
	fmt.Printf("%v successfully added and followed %v\n", feedRows.UserName, feedRows.FeedName)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(feeds) == 0 {
		fmt.Println("no feeds")
		os.Exit(1)
	}
	for _, feed := range feeds {
		fmt.Printf("* %v, %v, %v\n", feed.Name, feed.Url, feed.Name_2)
	}
	os.Exit(0)
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Please input the URL of the feed you want to follow.")
		os.Exit(1)
	}
	lookup, err := s.Db.FeedLookup(context.Background(), cmd.Args[0])
	if lookup == (database.FeedLookupRow{}) {
		fmt.Println("Feed not registered, add feed before you follow.")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("An Error occured when looking up the feed.", err)
		os.Exit(1)
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    lookup.ID,
	}
	feedRows, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		fmt.Println("Error when following feed", err)
		os.Exit(1)
	}
	fmt.Printf("%v successfully followed %v", feedRows.UserName, feedRows.FeedName)
	os.Exit(0)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	feedList, err := s.Db.CreateFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Println("an error occured with getting following list: ", err)
		os.Exit(1)
	}
	for _, f := range feedList {
		fmt.Printf("* %v - %v\n", f.FeedName, f.UserName)
	}
	os.Exit(0)
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Please input the URL of the feed you want to unfollow.")
		os.Exit(1)
	}
	lookup, err := s.Db.FeedLookup(context.Background(), cmd.Args[0])
	if lookup == (database.FeedLookupRow{}) {
		fmt.Println("Feed not registered or followed.")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("An Error occured when looking up the feed.", err)
		os.Exit(1)
	}
	feedList, err := s.Db.CreateFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Println("an error occured with getting following list: ", err)
		os.Exit(1)
	}
	if len(feedList) == 0 {
		fmt.Println("You are not following anyone.")
		os.Exit(1)
	}
	for _, f := range feedList {
		if lookup.ID == f.FeedID {
			params := database.RemoveFeedFollowParams{
				FeedID: lookup.ID,
				UserID: user.ID,
			}
			if err = s.Db.RemoveFeedFollow(context.Background(), params); err != nil {
				fmt.Println("An error occurred removing feed from database: ", err)
				os.Exit(1)
			}
			fmt.Println("Removed feed successfully.")
			os.Exit(0)
			return nil
		}
	}
	fmt.Println("You are not following that feed.")
	os.Exit(1)
	return fmt.Errorf("user not following given url\n")

}

func HandlerBrowse(s *State, cmd Command) error {
	limit := 2
	if len(cmd.Args) > 0 {
		temp, err := strconv.Atoi(cmd.Args[0])
		limit = temp
		if err != nil {
			return err
		}
	}
	posts, err := s.Db.GetPosts(context.Background(), int32(limit))
	if err != nil {
		fmt.Println("Error when getting posts: ", err)
	}
	for _, post := range posts {
		fmt.Printf("Title: %v\n", post.Title)
		fmt.Printf("Description: %v\n", post.Description)
		fmt.Printf("Published at: %v\n", post.PublishedAt)
		fmt.Printf("Link for more details: %v\n", post.Url)
	}
	return nil
}

func MiddlewareLoggedIn(
	handler func(s *State, cmd Command, user database.User) error,
) func(*State, Command) error {

	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("error when getting next feed to fetch: ", err)
		os.Exit(1)
	}
	params := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        feed.ID,
	}
	if err = s.Db.MarkFeedFetched(context.Background(), params); err != nil {
		fmt.Println("error when marking feed as fetched: ", err)
	}
	RSS, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Println("Error with FetchFeed func: ", err)
		return err
	}
	for _, item := range RSS.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return err
		}
		PostParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		}
		s.Db.CreatePost(context.Background(), PostParams)
	}
	return nil
}
