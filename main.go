package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// metricValue holds the current value of our custom metric.
// We use a mutex to ensure thread-safe access.
var (
	metricValue float64
	mu          sync.Mutex
)

// Create a new Prometheus gauge metric.
// A gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
var customMetric = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "custom_metric_value",
	Help: "The current value of the custom settable metric.",
})

// viewHandler serves the main page, which displays the current metric value
// and provides a form to set a new one.
func viewHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Prometheus Custom Metric</title>
        <style>
            body { font-family: Arial, sans-serif; margin: 40px; }
            .container { max-width: 600px; margin: auto; padding: 20px; border: 1px solid #ddd; border-radius: 8px; }
            h1 { color: #333; }
            p { font-size: 1.2em; }
            form { margin-top: 20px; }
            input[type="number"] { padding: 8px; width: 100px; }
            input[type="submit"] { padding: 8px 15px; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Custom Metric Value</h1>
            <p>Current Value: <strong>%v</strong></p>
            <form action="/set" method="post">
                <label for="value">Set New Value:</label>
                <input type="number" id="value" name="value" required>
                <input type="submit" value="Set">
            </form>
        </div>
    </body>
    </html>
    `
	fmt.Fprintf(w, html, metricValue)
}

// setHandler handles the API request to update the metric's value.
// It expects a POST request with a 'value' form field.
func setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	valueStr := r.FormValue("value")
	if valueStr == "" {
		http.Error(w, "Missing 'value' form field", http.StatusBadRequest)
		return
	}

	newVal, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		http.Error(w, "Invalid value. Please provide an integer.", http.StatusBadRequest)
		return
	}

	mu.Lock()
	metricValue = newVal
	mu.Unlock()

	// Update the Prometheus gauge.
	customMetric.Set(metricValue)

	log.Printf("Metric value set to: %v\n", metricValue)

	// Redirect back to the main page to show the new value.
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	// Initialize the metric with a default value of 0.
	customMetric.Set(0)

	// Handler for the main page.
	http.HandleFunc("/", viewHandler)

	// Handler for setting the metric value.
	http.HandleFunc("/set", setHandler)

	// Expose the registered metrics via the /metrics endpoint for Prometheus.
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

