# Triple Extraction for Ocean Science Applications

Ocean science is an excellent domain for applying your triple extraction and rule processing approach. The field generates diverse, heterogeneous data from multiple sources that need to be semantically integrated for knowledge discovery.

## Key Applications in Ocean Science

1. **Oceanographic Data Integration**
   - Extract semantic relationships from diverse instrument outputs (buoys, gliders, satellites)
   - Integrate multi-modal observations (temperature, salinity, current speed) into knowledge graphs
   - Enable cross-dataset queries without requiring uniform data formats

2. **Marine Ecosystem Monitoring**
   - Process observational data about species interactions
   - Extract relationships between environmental conditions and biological responses
   - Build knowledge bases connecting observed phenomena to underlying causes

3. **Climate Change Research**
   - Extract temporal relationships between ocean measurements
   - Identify correlations across different timescales and regions
   - Support attribution of changes to specific factors

## Example Implementation

```go
// OceanScienceExtractor uses the triple extraction framework for ocean data
package oceansci

import (
  "context"
  "fmt"
  "time"

  "github.com/google/uuid"
  "github.com/witfoo/krapht/triple/extraction"
)

// BuoyReading represents data from an oceanographic buoy
type BuoyReading struct {
  ID           string    `json:"id"`
  StationID    string    `json:"station_id"`
  Timestamp    time.Time `json:"timestamp"`
  Latitude     float64   `json:"latitude"`
  Longitude    float64   `json:"longitude"`
  Temperature  float64   `json:"temperature"`
  Salinity     float64   `json:"salinity"`
  WaveHeight   float64   `json:"wave_height"`
  WindSpeed    float64   `json:"wind_speed"`
  WindDirection int      `json:"wind_direction"`
  Pressure     float64   `json:"pressure"`
}

// CreateDefaultOceanRules creates basic extraction rules for ocean data
func CreateDefaultOceanRules(ctx context.Context, store *extraction.RuleStore) error {
  // Add buoy measurement rules
  buoyRule := &extraction.ExtractionRule{
    ID:         "buoy-measurements",
    Name:       "Buoy Measurement Relationships",
    Description: "Extracts relationships between buoys and measurements",
    SourceType: "buoy_reading",
    Condition: extraction.Condition{
      Field:    "StationID", 
      Operator: "neq",
      Value:    "",
    },
    Extractions: []extraction.ExtractionPattern{
      {
        // Location relationship
        Subject: extraction.FieldSelector{
          Type:  "field",
          Value: "StationID",
        },
        Predicate: "has_location",
        Object: extraction.FieldSelector{
          Type:  "template",
          Value: "{{.Latitude}},{{.Longitude}}",
        },
      },
      {
        // Temperature measurement
        Subject: extraction.FieldSelector{
          Type:  "field",
          Value: "StationID",
        },
        Predicate: "measures_temperature",
        Object: extraction.FieldSelector{
          Type:  "field",
          Value: "Temperature",
        },
      },
      {
        // Salinity measurement
        Subject: extraction.FieldSelector{
          Type:  "field",
          Value: "StationID",
        },
        Predicate: "measures_salinity",
        Object: extraction.FieldSelector{
          Type:  "field",
          Value: "Salinity",
        },
      },
      {
        // Wave properties
        Subject: extraction.FieldSelector{
          Type:  "field",
          Value: "StationID",
        },
        Predicate: "observes_wave_height",
        Object: extraction.FieldSelector{
          Type:  "field",
          Value: "WaveHeight",
        },
      },
    },
    Priority: 100,
  }

  if err := store.AddRule(ctx, buoyRule); err != nil {
    return fmt.Errorf("adding buoy rule: %w", err)
  }

  // Add weather condition rule that extracts weather-related triples
  weatherRule := &extraction.ExtractionRule{
    ID:          "weather-conditions",
    Name:        "Weather Condition Relationships",
    Description: "Extracts relationships about weather conditions",
    SourceType:  "buoy_reading",
    Condition:   extraction.Condition{
      Field:    "WindSpeed",
      Operator: "gte",
      Value:    0,
    },
    Extractions: []extraction.ExtractionPattern{
      {
        Subject: extraction.FieldSelector{
          Type:  "field",
          Value: "StationID",
        },
        Predicate: "measures_wind_speed",
        Object: extraction.FieldSelector{
          Type:  "field",
          Value: "WindSpeed",
        },
      },
      {
        Subject: extraction.FieldSelector{
          Type:  "field",
          Value: "StationID", 
        },
        Predicate: "measures_pressure",
        Object: extraction.FieldSelector{
          Type:  "field",
          Value: "Pressure",
        },
      },
    },
    Priority: 90,
  }

  if err := store.AddRule(ctx, weatherRule); err != nil {
    return fmt.Errorf("adding weather rule: %w", err)
  }

  // Rule for storm detection
  stormRule := &extraction.ExtractionRule{
    ID:          "storm-detection",
    Name:        "Storm Condition Detection",
    Description: "Detects when measurements indicate storm conditions",
    SourceType:  "buoy_reading",
    Condition: extraction.Condition{
      LogicalOp: "and",
      Children: []extraction.Condition{
        {
          Field:    "WindSpeed",
          Operator: "gt",
          Value:    35.0, // Storm force winds
        },
        {
          Field:    "WaveHeight",
          Operator: "gt",
          Value:    3.5, // High waves
        },
      },
    },
    Extractions: []extraction.ExtractionPattern{
      {
        Subject: extraction.FieldSelector{
          Type:  "field",
          Value: "StationID",
        },
        Predicate: "indicates_storm_condition",
        Object: extraction.FieldSelector{
          Type:  "static",
          Value: "true",
        },
      },
    },
    Priority: 100,
  }

  if err := store.AddRule(ctx, stormRule); err != nil {
    return fmt.Errorf("adding storm rule: %w", err)
  }

  return nil
}

// ProcessOceanographicData processes oceanographic data and extracts triples
func ProcessOceanographicData(ctx context.Context, processor *extraction.ArtifactProcessor, data []BuoyReading) error {
  for _, reading := range data {
    select {
    case <-ctx.Done():
      return ctx.Err()
    default:
      if err := processor.ProcessArtifact(ctx, reading, "buoy_reading"); err != nil {
        return fmt.Errorf("processing buoy reading %s: %w", reading.ID, err)
      }
    }
  }
  return nil
}
```

## Benefits for Ocean Science

1. **Integration of Heterogeneous Data**
   - Ocean science data comes from diverse sources (satellites, buoys, ships, underwater vehicles)
   - Triple extraction normalizes these into a consistent semantic framework

2. **Complex Pattern Recognition**
   - Identify oceanographic phenomena spanning multiple measurements
   - Detect correlations between physical, chemical, and biological variables

3. **Knowledge Accumulation**
   - Build knowledge graphs of ocean phenomena over time
   - Preserve semantic relationships across multiple observation campaigns

4. **Cross-Domain Integration**
   - Connect ocean data with climate models, weather systems, and marine biology
   - Enable interdisciplinary queries across traditionally siloed research areas

This approach allows oceanographers to focus on scientific insights rather than data wrangling, while providing a flexible framework that can adapt to new instruments and observation types without requiring code changes.
