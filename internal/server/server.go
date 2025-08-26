package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/blesniewski/knm/internal/models"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer         *http.Server
	exchangeRateClient ExchangeRateClient
	cryptoClient       CryptoConversionClient
}

type ExchangeRateClient interface {
	GetRatesForCurrencies(ctx context.Context, currencies []string) ([]models.CurrencyPair, error)
}

type CryptoConversionClient interface {
	GetConversionRate(from, to string, amount float64) (models.CryptoPair, error)
}

func New(exchangeRateClient ExchangeRateClient, cryptoClient CryptoConversionClient) *Server {
	s := &Server{
		exchangeRateClient: exchangeRateClient,
		cryptoClient:       cryptoClient,
	}
	s.registerNewRoutes()
	return s
}

func (s *Server) Run(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) registerNewRoutes() {
	router := gin.Default()
	router.GET("/rates", s.handleGetRates)
	router.GET("/exchange", s.handleGetExchange)
	s.httpServer = &http.Server{
		Handler: router,
	}
}

func (s *Server) handleGetRates(c *gin.Context) {
	currencies := c.Query("currencies")
	currencyList := strings.Split(currencies, ",")
	if len(currencyList) < 2 {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	rates, err := s.exchangeRateClient.GetRatesForCurrencies(c.Request.Context(), currencyList)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	c.JSON(http.StatusOK, rates)
}

func (s *Server) handleGetExchange(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")
	if from == "" || to == "" || amountStr == "" {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	exchange, err := s.cryptoClient.GetConversionRate(from, to, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	c.JSON(http.StatusOK, exchange)
}
