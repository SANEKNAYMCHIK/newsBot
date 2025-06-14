package main

import (
	"context"
	"log"
	"net"

	"github.com/SANEKNAYMCHIK/newsBot/internal/app/parser"
	"github.com/SANEKNAYMCHIK/newsBot/proto/news"
	"google.golang.org/grpc"
)

type server struct {
	news.UnimplementedNewsParserServer
}

func (s *server) GetNews(ctx context.Context, req *news.NewsRequest) (*news.NewsResponse, error) {
	result := parser.ParseAllNews(req.Sources)
	data := make(map[string]*news.NewsItemList)
	for name, newsList := range result {
		itemList := &news.NewsItemList{
			Items: make([]*news.NewsItem, 0, len(newsList)),
		}
		for _, n := range newsList {
			itemList.Items = append(itemList.Items, &news.NewsItem{
				Title:       n.Title,
				Link:        n.Link,
				Date:        n.Date.Format("2006-01-02 15:04:05"),
				Description: n.Description,
				Website:     n.Website,
			})
		}
		data[name] = itemList
	}
	// for _, newsList := range result {
	// 	for _, n := range newsList {
	// 		items = append(items, &news.NewsItem{
	// 			Title:       n.Title,
	// 			Link:        n.Link,
	// 			Date:        n.Date.Format("2006-01-02 15:04:05"),
	// 			Description: n.Description,
	// 			Website:     n.Website,
	// 		})
	// 	}
	// }
	// return &news.NewsResponse{Items: items}, nil
	return &news.NewsResponse{Data: data}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	news.RegisterNewsParserServer(s, &server{})
	log.Println("newsParser gRPC server started on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
