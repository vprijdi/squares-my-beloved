import React, { useState, useMemo } from 'react';

const ScoreTrendGraph = () => {
  const [hoveredPoint, setHoveredPoint] = useState(null);

  const data = useMemo(() => {
    const points = [];
    const today = new Date();

    // Find the most recent Monday
    const daysFromMonday = (today.getDay() + 6) % 7;
    const lastMonday = new Date(today);
    lastMonday.setDate(today.getDate() - daysFromMonday);

    // Predefined score progression that ends up higher
    const scoreChanges = [
      0, -12, 8, -23, 15, 34, -18, 22, -8, 16, -31, 28, 12, -19, 37, -14, 25, 9,
      -27, 33, -5, 18, -22, 41, 7, -16, 29, -11, 24, 13, -35, 42, -7, 20, -13,
      31, 6, -24, 38, -9, 17, -26, 45, 2, -18, 23, -4, 36, 14, -21, 39, 28,
    ];

    let score = 52; // Starting score

    // Generate 52 data points going backwards
    for (let i = 51; i >= 0; i--) {
      const date = new Date(lastMonday);
      date.setDate(lastMonday.getDate() - i * 7);

      // Apply predefined score change
      if (i < 51) {
        score += scoreChanges[51 - i];
      }

      points.push({
        date: date.toISOString().split('T')[0],
        score: score,
        week: 52 - i,
      });
    }

    return points;
  }, []);

  // Calculate graph dimensions and scaling
  const width = 1000;
  const height = 400;
  const padding = 60;
  const graphWidth = width - 2 * padding;
  const graphHeight = height - 2 * padding;

  const minScore = Math.min(...data.map((d) => d.score));
  const maxScore = Math.max(...data.map((d) => d.score));
  const scoreRange = maxScore - minScore || 1;

  // Create path for the line
  const pathData = data
    .map((point, index) => {
      const x = padding + (index / (data.length - 1)) * graphWidth;
      const y = padding + ((maxScore - point.score) / scoreRange) * graphHeight;
      return `${index === 0 ? 'M' : 'L'} ${x} ${y}`;
    })
    .join(' ');

  // Create points for hover detection
  const points = data.map((point, index) => ({
    ...point,
    x: padding + (index / (data.length - 1)) * graphWidth,
    y: padding + ((maxScore - point.score) / scoreRange) * graphHeight,
  }));

  const formatDate = (dateStr) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  return (
    <div className="score-graph-container">
      <style jsx>{`
        .score-graph-container {
          background-color: #1a1a1a;
          color: white;
          padding: 20px;
          font-family:
            -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
          min-height: 100vh;
          display: flex;
          flex-direction: column;
          align-items: center;
        }

        .graph-title {
          font-size: 1.25rem;
          font-weight: bold;
          margin-bottom: 20px;
          text-align: center;
        }

        .graph-svg {
          background-color: #1a1a1a;
          border: 1px solid #333;
          border-radius: 8px;
        }

        .graph-line {
          fill: none;
          stroke: white;
          stroke-width: 2;
        }

        .graph-point {
          fill: white;
          cursor: pointer;
          transition: r 0.2s ease;
        }

        .graph-point:hover {
          r: 6;
        }

        .axis-line {
          stroke: #666;
          stroke-width: 1;
        }

        .axis-text {
          fill: white;
          font-size: 12px;
          text-anchor: middle;
        }

        .axis-text.y-axis {
          text-anchor: end;
        }

        .tooltip {
          position: absolute;
          background-color: rgba(0, 0, 0, 0.9);
          color: white;
          padding: 8px 12px;
          border-radius: 4px;
          font-size: 14px;
          pointer-events: none;
          border: 1px solid #666;
          z-index: 10;
        }

        .stats {
          margin-top: 20px;
          display: flex;
          gap: 40px;
          font-size: 14px;
        }

        .stat-item {
          text-align: center;
        }

        .stat-value {
          font-size: 18px;
          font-weight: bold;
          color: #fff;
        }

        .stat-label {
          color: #ccc;
          margin-top: 4px;
        }
      `}</style>

      <h1 className="graph-title">Weekly Score Trend - Last 52 Mondays</h1>

      <div style={{ position: 'relative' }}>
        <svg
          className="graph-svg"
          width={width}
          height={height}
          onMouseLeave={() => setHoveredPoint(null)}
        >
          {/* Y-axis */}
          <line
            className="axis-line"
            x1={padding}
            y1={padding}
            x2={padding}
            y2={height - padding}
          />

          {/* X-axis */}
          <line
            className="axis-line"
            x1={padding}
            y1={height - padding}
            x2={width - padding}
            y2={height - padding}
          />

          {/* Y-axis labels */}
          {[0, 0.25, 0.5, 0.75, 1].map((ratio) => {
            const score = Math.round(minScore + (maxScore - minScore) * ratio);
            const y = height - padding - ratio * graphHeight;
            return (
              <g key={ratio}>
                <line
                  className="axis-line"
                  x1={padding - 5}
                  y1={y}
                  x2={padding}
                  y2={y}
                />
                <text className="axis-text y-axis" x={padding - 10} y={y + 4}>
                  {score}
                </text>
              </g>
            );
          })}

          {/* X-axis labels (every 8 weeks) */}
          {data
            .filter((_, index) => index % 8 === 0)
            .map((point, index) => {
              const x =
                padding +
                (data.indexOf(point) / (data.length - 1)) * graphWidth;
              return (
                <g key={index}>
                  <line
                    className="axis-line"
                    x1={x}
                    y1={height - padding}
                    x2={x}
                    y2={height - padding + 5}
                  />
                  <text className="axis-text" x={x} y={height - padding + 20}>
                    {formatDate(point.date).split(',')[0]}
                  </text>
                </g>
              );
            })}

          {/* Graph line */}
          <path className="graph-line" d={pathData} />

          {/* Data points */}
          {points.map((point, index) => (
            <circle
              key={index}
              className="graph-point"
              cx={point.x}
              cy={point.y}
              r={hoveredPoint === index ? 6 : 4}
              onMouseEnter={() => setHoveredPoint(index)}
            />
          ))}
        </svg>

        {/* Tooltip */}
        {hoveredPoint !== null && (
          <div
            className="tooltip"
            style={{
              left: points[hoveredPoint].x + 10,
              top: points[hoveredPoint].y - 40,
            }}
          >
            <div>
              <strong>Date:</strong> {formatDate(data[hoveredPoint].date)}
            </div>
            <div>
              <strong>Score:</strong> {data[hoveredPoint].score}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ScoreTrendGraph;
