/*
 * Copyright (C) 2022 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package node

import (
	"github.com/pkg/errors"

	"github.com/mysteriumnetwork/node/identity"
)

// MonitoringAgentStatuses a object represent a [service_type][status]amount of statuses for each service type.
type MonitoringAgentStatuses map[string]map[string]int

// ProviderStatuses should return provider statuses from monitoring agent
type ProviderStatuses func(providerID string) (MonitoringAgentStatuses, error)

// ProviderSessionsList should return provider sessions list
type ProviderSessionsList func(id identity.Identity, rangeTime string) ([]SessionItem, error)

// ProviderTransferredData should return total traffic served by the provider during a period of time
type ProviderTransferredData func(id identity.Identity, rangeTime string) (TransferredData, error)

// ProviderSessionsCount should return provider sessions count
type ProviderSessionsCount func(id identity.Identity, rangeTime string) (SessionsCount, error)

// ProviderConsumersCount should return unique consumers count
type ProviderConsumersCount func(id identity.Identity, rangeTime string) (ConsumersCount, error)

// StatsTracker tracks metrics for service
type StatsTracker struct {
	providerStatuses        ProviderStatuses
	providerSessionsList    ProviderSessionsList
	providerTransferredData ProviderTransferredData
	providerSessionsCount   ProviderSessionsCount
	providerConsumersCount  ProviderConsumersCount
	currentIdentity         currentIdentity
}

// NewNodeStatsTracker constructor
func NewNodeStatsTracker(
	providerStatuses ProviderStatuses,
	providerSessions ProviderSessionsList,
	providerTransferredData ProviderTransferredData,
	providerSessionsCount ProviderSessionsCount,
	providerConsumersCount ProviderConsumersCount,
	currentIdentity currentIdentity,
) *StatsTracker {
	mat := &StatsTracker{
		providerStatuses:        providerStatuses,
		providerSessionsList:    providerSessions,
		providerTransferredData: providerTransferredData,
		providerSessionsCount:   providerSessionsCount,
		providerConsumersCount:  providerConsumersCount,
		currentIdentity:         currentIdentity,
	}

	return mat
}

// Statuses retrieves and resolved monitoring status from quality oracle
func (m *StatsTracker) Statuses() (MonitoringAgentStatuses, error) {
	id, ok := m.currentIdentity.GetUnlockedIdentity()
	if ok {
		return m.providerStatuses(id.Address)
	}

	return MonitoringAgentStatuses{}, errors.New("identity not found")
}

// SessionItem represents information about session monitoring metrics.
type SessionItem struct {
	ID              string `json:"id"`
	ConsumerCountry string `json:"consumer_country"`
	ServiceType     string `json:"service_type"`
	Duration        int64  `json:"duration"`
	StartedAt       int64  `json:"started_at"`
	Earning         string `json:"earning"`
	Transferred     int64  `json:"transferred"`
}

// TransferredData represent information about total traffic served by the provider during a period of time
type TransferredData struct {
	Bytes int `json:"transferred_data_bytes"`
}

// SessionsCount represent a information about number of sessions during a period of time
type SessionsCount struct {
	Count int `json:"count"`
}

// ConsumersCount represent a information about number of consumers served during a period of time
type ConsumersCount struct {
	Count int `json:"count"`
}

// Sessions retrieves and resolved monitoring status from quality oracle
func (m *StatsTracker) Sessions(rangeTime string) ([]SessionItem, error) {
	id, ok := m.currentIdentity.GetUnlockedIdentity()
	if ok {
		return m.providerSessionsList(id, rangeTime)
	}

	return []SessionItem{}, errors.New("identity not found")
}

// TransferredData retrieves and resolved total traffic served by the provider
func (m *StatsTracker) TransferredData(rangeTime string) (TransferredData, error) {
	id, ok := m.currentIdentity.GetUnlockedIdentity()
	if ok {
		return m.providerTransferredData(id, rangeTime)
	}

	return TransferredData{}, errors.New("identity not found")
}

// SessionsCount retrieves and resolved numbers of sessions
func (m *StatsTracker) SessionsCount(rangeTime string) (SessionsCount, error) {
	id, ok := m.currentIdentity.GetUnlockedIdentity()
	if ok {
		return m.providerSessionsCount(id, rangeTime)
	}

	return SessionsCount{}, errors.New("identity not found")
}

// ConsumersCount retrieves and resolved numbers of consumers server during period of time
func (m *StatsTracker) ConsumersCount(rangeTime string) (ConsumersCount, error) {
	id, ok := m.currentIdentity.GetUnlockedIdentity()
	if ok {
		return m.providerConsumersCount(id, rangeTime)
	}

	return ConsumersCount{}, errors.New("identity not found")
}
