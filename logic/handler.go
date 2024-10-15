package logic

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/oklog/ulid"
)

func GenerateQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pgTrxMethodQueries := make([]string, 0)
	pgTrxMethodIDs := make([]string, 0)

	pgChanSettingQueries := make([]string, 0)
	pgChanSettingIDs := make([]string, 0)

	// loop request data
	for _, data := range req.Data {
		// generate pg trx method query
		pgTrxMethodQuery, pgTrxMethodID := genPgTrxMethodQuery(data.TMMID, req.PgAccountID, data.ChannelCode)
		pgTrxMethodIDs = append(pgTrxMethodIDs, fmt.Sprintf(`'%s',`, pgTrxMethodID))
		pgTrxMethodQueries = append(pgTrxMethodQueries, pgTrxMethodQuery)

		// generate pg chan setting query
		pgChanSettingQuery, pgChanSettingID := genPgChanSettingQuery(pgTrxMethodID)
		pgChanSettingIDs = append(pgChanSettingIDs, fmt.Sprintf(`'%s',`, pgChanSettingID))
		pgChanSettingQueries = append(pgChanSettingQueries, pgChanSettingQuery)
	}

	// Write arrays to file
	err = writeArraysToFile(pgTrxMethodQueries, pgChanSettingQueries, pgChanSettingIDs, pgTrxMethodIDs)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	// Respond with the received data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// ('ptm-01HRGLCNPGF68VRNF2HQS0FVZR', 'pga-01J52NFCGJ529TJ84PTT9SE4JQ', 'ACTIVE', 'tmm-01HYJZP35Z62CY4QB1JQMN27M2', '097'),
func genPgTrxMethodQuery(tmmID, pgAccountID, channelCode string) (string, string) {
	pgTrxMethodID, _ := GeneratePK("ptm")
	return fmt.Sprintf(`('%s', '%s', 'ACTIVE', '%s', '%s'),`, pgTrxMethodID, pgAccountID, tmmID, channelCode), pgTrxMethodID
}

// ('ppcs-01HRGKTQ2MKPP4EZQXA9SXJDPB', 'ptm-01HRGKP1XGQJ0TKVBWXM8Y1Q0N', 'payouts', 0, 0, 'NO_COST', 0, 0, 0, 0),
func genPgChanSettingQuery(pgTrxMethodID string) (string, string) {
	pgChanSettingID, _ := GeneratePK("ppcs")
	return fmt.Sprintf(`('%s', '%s', 'payouts', 0, 0, 'NO_COST', 0, 0, 0, 0),`, pgChanSettingID, pgTrxMethodID), pgChanSettingID
}

func GeneratePK(prefix string) (id string, err error) {
	// Create a source of randomness from crypto/rand.
	entropy := ulid.Monotonic(rand.Reader, 0)
	ms := ulid.Timestamp(time.Now())
	ULID, err := ulid.New(ms, entropy)
	if err != nil {
		return
	}

	id = fmt.Sprintf("%s-%s", prefix, ULID.String())
	return id, nil
}

func writeArraysToFile(pgTrxMethodQueries, pgChanSettingQueries, pgChanSettingIDs, pgTrxMethodIDs []string) error {
	file, err := os.Create("output.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write pgTrxMethodQueries
	_, err = file.WriteString("pgTrxMethodQueries:\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString(strings.Join(pgTrxMethodQueries, "\n") + "\n\n")
	if err != nil {
		return err
	}

	// Write pgChanSettingQueries
	_, err = file.WriteString("pgChanSettingQueries:\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString(strings.Join(pgChanSettingQueries, "\n") + "\n\n")
	if err != nil {
		return err
	}

	// Write pgChanSettingIDs
	_, err = file.WriteString("pgChanSettingIDs:\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString(strings.Join(pgChanSettingIDs, "\n") + "\n\n")
	if err != nil {
		return err
	}

	// Write pgTrxMethodIDs
	_, err = file.WriteString("pgTrxMethodIDs:\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString(strings.Join(pgTrxMethodIDs, "\n") + "\n")
	if err != nil {
		return err
	}

	return nil
}
