package api

import (
	"fmt"
	"net/http"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

// listUserStatsHandler godoc
//
//	@Summary		Get user statistics
//	@Description	Retrieve various statistics for the authenticated user. Categories must be comma-separated if providing multiple values.
//	@Tags			statistics
//	@Accept			json
//	@Produce		json
//	@Param			categories	query	string	false	"Comma-separated categories (all_time,tasks,time,averages,streaks)"	default(all_time,tasks,time,averages,streaks)
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]interface{}	"Statistics returned successfully"
//	@Failure		400	{object}	map[string]string		"Bad request: bad category value or malformed input"
//	@Failure		401	{object}	map[string]string		"Unauthorized: missing or invalid authentication"
//	@Failure		404	{object}	map[string]string		"Not Found: user statistics not available"
//	@Failure		500	{object}	map[string]string		"Internal Server Error"
//	@Router			/stats [get]
func (app *Application) listUserStatsHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	sf := &store.StatisticsFilters{
		Categories: []string{},
	}

	sf, err := sf.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(sf); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	stats, err := app.store.Statistics.GetByUserID(ctx, user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	response := make(map[string]interface{})
	data := make(map[string]interface{})

	for _, cat := range sf.Categories {
		switch cat {
		case "all_time":
			data["all_time"] = map[string]interface{}{
				"highest_daily_score":      stats.HighestDailyScore,
				"highest_daily_score_date": stats.HighestDailyScoreDate,
				"lowest_daily_score":       stats.LowestDailyScore,
				"lowest_daily_score_date":  stats.LowestDailyScoreDate,
			}
		case "tasks":
			data["tasks"] = map[string]interface{}{
				"total_completed":    stats.TotalTasksCompleted,
				"tier0_completed":    stats.TotalTier0Completed,
				"tier1_completed":    stats.TotalTier1Completed,
				"tier2_completed":    stats.TotalTier2Completed,
				"avg_per_active_day": stats.AvgTasksPerActiveDay,
			}
		case "time":
			data["time"] = map[string]interface{}{
				"total_active_days":        stats.TotalActiveDays,
				"total_time_spent_minutes": stats.TotalTimeSpentMinutes,
			}
		case "averages":
			data["averages"] = map[string]interface{}{
				"daily_score":   stats.AvgDailyScore,
				"weekly_score":  stats.AvgWeeklyScore,
				"monthly_score": stats.AvgMonthlyScore,
			}
		case "streaks":
			data["streaks"] = map[string]interface{}{
				"current_days":       stats.CurrentStreakDays,
				"longest_days":       stats.LongestStreakDays,
				"longest_start_date": stats.LongestStreakStartDate,
				"longest_end_date":   stats.LongestStreakEndDate,
			}
		default:
			app.badRequestResponse(w, r, fmt.Errorf("invalid category: %s", cat))
			return
		}
	}

	response["data"] = data
	response["meta"] = map[string]interface{}{
		"categories": sf.Categories,
		"updated_at": stats.UpdatedAt,
	}

	if err := writeJSON(w, http.StatusOK, response); err != nil {
		app.internalServerError(w, r, err)
	}
}
