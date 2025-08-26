package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/blesniewski/knm/internal/models"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router             *gin.Engine
	exchangeRateClient ExchangeRateClient
	cryptoClient       CryptoConversionClient
}

type ExchangeRateClient interface {
	GetRatesForCurrencies(currencies []string) ([]models.CurrencyPair, error)
}

type CryptoConversionClient interface {
	GetConversionRate(from, to string, amount float64) (models.CryptoPair, error)
}

func NewServer(exchangeRateClient ExchangeRateClient, cryptoClient CryptoConversionClient) *Server {
	s := &Server{
		exchangeRateClient: exchangeRateClient,
		cryptoClient:       cryptoClient,
	}
	s.registerNewRoutes()
	return s
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) registerNewRoutes() {
	router := gin.Default()
	router.GET("/rates", s.handleGetRates)
	router.GET("/exchange", s.handleGetExchange)
	s.router = router
}

func (s *Server) handleGetRates(c *gin.Context) {
	currencies := c.Query("currencies")
	currencyList := strings.Split(currencies, ",")
	if len(currencyList) < 2 {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	rates, err := s.exchangeRateClient.GetRatesForCurrencies(currencyList)
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
