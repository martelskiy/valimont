package validator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"io"

	"github.com/martelskiy/valimont/internal/attestation"
	"golang.org/x/time/rate"
)

const (
	scheme  = "https"
	baseURL = "holesky.beaconcha.in"
)

type Validator struct {
	client        *http.Client
	ratelimiter   *rate.Limiter
	validatorIndx string
	validatorsNr  uint16
}

func New(validatorIndex []uint32, rateLimitPerMinute uint8) *Validator {
	rl := rate.NewLimiter(rate.Every(time.Minute), int(rateLimitPerMinute))

	strNumbers := make([]string, len(validatorIndex))
	for i, num := range validatorIndex {
		strNumbers[i] = strconv.Itoa(int(num))
	}

	return &Validator{
		client:        http.DefaultClient,
		ratelimiter:   rl,
		validatorIndx: strings.Join(strNumbers, ","),
		validatorsNr:  uint16(len(validatorIndex)),
	}
}

func (v *Validator) GetAttestations(ctx context.Context) ([]attestation.Model, error) {
	var attestations []attestation.Model
	r, err := v.getResponse(ctx)
	if err != nil {
		return attestations, errors.Join(err, errors.New("error waiting for rate limiter"))
	}

	var maxTargetEpoch uint64
	epochData := make(map[uint64][]attestation.Model)
	for _, att := range r.Data {
		epochData[uint64(att.Epoch)] = append(epochData[uint64(att.Epoch)], attestation.Model{
			Epoch:             att.Epoch,
			Status:            att.Status,
			InclusionDistance: att.InclusionSlot - att.AttesterSlot,
		})

		if maxTargetEpoch < uint64(att.Epoch) {
			maxTargetEpoch = uint64(att.Epoch)
		}
	}

	return epochData[maxTargetEpoch-2], nil
}

func (v *Validator) getResponse(ctx context.Context) (response, error) {
	var r response
	if err := v.ratelimiter.Wait(ctx); err != nil {
		return r, errors.Join(err, errors.New("error waiting for rate limiter"))
	}

	resp, err := v.client.Get(fmt.Sprintf("%s://%s/api/v1/validator/%s/attestations", scheme, baseURL, v.validatorIndx))
	if err != nil {
		return r, errors.Join(err, errors.New("error from http client"))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return r, errors.Join(err, errors.New("failed to read response body"))
	}

	if err = json.Unmarshal(body, &r); err != nil {
		return r, errors.Join(err, errors.New("failed to unmarshal HTTP response"))
	}
	if r.Status != "OK" {
		return r, errors.Join(err, errors.New("response status was not 'OK'"))
	}
	return r, nil
}

func (v *Validator) ValidatorsNr() uint16 {
	return v.validatorsNr
}
