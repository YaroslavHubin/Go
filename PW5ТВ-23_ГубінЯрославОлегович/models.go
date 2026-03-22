package main

type Station struct {
    ID       int
    Name     string
    Power    int
    Slots    int
    Location string
}

type Session struct {
    ID        int
    StationID int
    Duration  int
    Energy    int
    StartedAt string
}

type Stat struct {
    ID            int
    Name          string
    Location      string
    AvgDuration   float64
    MaxEnergy     int
    TotalSessions int
}
