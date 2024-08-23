export function calculatePercent(
  numerator: number,
  denominator: number,
): number {
  if (denominator === 0) return 0;
  return parseFloat(((numerator / denominator) * 100).toFixed(2));
}
