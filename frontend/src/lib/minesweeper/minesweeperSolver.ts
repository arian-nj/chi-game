const DIRS: [number, number][] = [
  [-1, -1], [-1, 0], [-1, 1],
  [0, -1],           [0, 1],
  [1, -1],  [1, 0],  [1, 1],
]

export function isInFirstClickSafeZone(
  r: number,
  c: number,
  safeRow: number,
  safeCol: number,
): boolean {
  return Math.abs(r - safeRow) <= 1 && Math.abs(c - safeCol) <= 1
}

export function createEmptyMines(height: number, width: number): boolean[][] {
  return Array.from({ length: height }, () =>
    Array.from({ length: width }, () => false),
  )
}

export function computeNeighborMines(
  mines: boolean[][],
  height: number,
  width: number,
): number[][] {
  const counts = Array.from({ length: height }, () =>
    Array.from({ length: width }, () => 0),
  )

  for (let row = 0; row < height; row++) {
    for (let col = 0; col < width; col++) {
      if (mines[row]![col]) continue

      let count = 0
      for (const [dr, dc] of DIRS) {
        if (mines[row + dr]?.[col + dc]) count++
      }
      counts[row]![col] = count
    }
  }

  return counts
}

export function placeRandomMines(
  mines: boolean[][],
  height: number,
  width: number,
  mineCount: number,
  safeRow: number,
  safeCol: number,
): void {
  for (let r = 0; r < height; r++) {
    for (let c = 0; c < width; c++) {
      mines[r]![c] = false
    }
  }

  const cells: [number, number][] = []
  for (let r = 0; r < height; r++) {
    for (let c = 0; c < width; c++) {
      if (isInFirstClickSafeZone(r, c, safeRow, safeCol)) continue
      cells.push([r, c])
    }
  }

  const minesToPlace = Math.min(mineCount, cells.length)
  for (let i = 0; i < minesToPlace; i++) {
    const j = i + Math.floor(Math.random() * (cells.length - i))
    ;[cells[i], cells[j]] = [cells[j]!, cells[i]!]
    const [row, col] = cells[i]!
    mines[row]![col] = true
  }
}

function getNeighbors(
  row: number,
  col: number,
  height: number,
  width: number,
): [number, number][] {
  const neighbors: [number, number][] = []
  for (const [dr, dc] of DIRS) {
    const r = row + dr
    const c = col + dc
    if (r >= 0 && c >= 0 && r < height && c < width) {
      neighbors.push([r, c])
    }
  }
  return neighbors
}

function floodRevealSolver(
  revealed: boolean[][],
  flagged: boolean[][],
  neighborMines: number[][],
  height: number,
  width: number,
  startRow: number,
  startCol: number,
): void {
  const queue: [number, number][] = [[startRow, startCol]]

  while (queue.length > 0) {
    const [row, col] = queue.shift()!
    if (row < 0 || col < 0 || row >= height || col >= width) continue
    if (revealed[row]![col] || flagged[row]![col]) continue

    revealed[row]![col] = true
    if (neighborMines[row]![col] !== 0) continue

    for (const [dr, dc] of DIRS) {
      queue.push([row + dr, col + dc])
    }
  }
}

type LogicStepResult = 'progress' | 'stuck' | 'invalid'

function applyLogicIteration(
  revealed: boolean[][],
  flagged: boolean[][],
  mines: boolean[][],
  neighborMines: number[][],
  height: number,
  width: number,
): LogicStepResult {
  let progress = false

  for (let row = 0; row < height; row++) {
    for (let col = 0; col < width; col++) {
      if (!revealed[row]![col]) continue

      const clue = neighborMines[row]![col]!
      const neighbors = getNeighbors(row, col, height, width)
      const hidden: [number, number][] = []
      let flagCount = 0

      for (const [r, c] of neighbors) {
        if (flagged[r]![c]) flagCount++
        else if (!revealed[r]![c]) hidden.push([r, c])
      }

      if (hidden.length === 0) continue

      if (flagCount === clue) {
        for (const [r, c] of hidden) {
          if (mines[r]![c]) return 'invalid'
          if (!revealed[r]![c]) {
            floodRevealSolver(revealed, flagged, neighborMines, height, width, r, c)
            progress = true
          }
        }
      }

      if (hidden.length === clue - flagCount) {
        for (const [r, c] of hidden) {
          if (!flagged[r]![c]) {
            flagged[r]![c] = true
            progress = true
          }
        }
      }
    }
  }

  return progress ? 'progress' : 'stuck'
}

export function isLogicSolvable(
  mines: boolean[][],
  height: number,
  width: number,
  startRow: number,
  startCol: number,
): boolean {
  if (mines[startRow]![startCol]) return false

  const neighborMines = computeNeighborMines(mines, height, width)
  const revealed = createEmptyMines(height, width)
  const flagged = createEmptyMines(height, width)

  floodRevealSolver(revealed, flagged, neighborMines, height, width, startRow, startCol)

  while (true) {
    const result = applyLogicIteration(revealed, flagged, mines, neighborMines, height, width)
    if (result === 'invalid') return false
    if (result === 'stuck') break
  }

  for (let row = 0; row < height; row++) {
    for (let col = 0; col < width; col++) {
      if (!mines[row]![col] && !revealed[row]![col]) return false
    }
  }

  return true
}

export function defaultMaxAttempts(width: number, height: number): number {
  return Math.max(100, width * height * 5)
}

export function generateSolvableMines(
  height: number,
  width: number,
  mineCount: number,
  safeRow: number,
  safeCol: number,
  maxAttempts = defaultMaxAttempts(width, height),
): { mines: boolean[][]; solvable: boolean } {
  const mines = createEmptyMines(height, width)

  for (let i = 0; i < maxAttempts; i++) {
    placeRandomMines(mines, height, width, mineCount, safeRow, safeCol)
    if (isLogicSolvable(mines, height, width, safeRow, safeCol)) {
      return { mines, solvable: true }
    }
  }

  placeRandomMines(mines, height, width, mineCount, safeRow, safeCol)
  return { mines, solvable: false }
}
