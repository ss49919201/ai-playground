/**
 * 二分探索の実装
 * ソート済み配列から対象要素を検索する
 *
 * 時間計算量: O(log n)
 * 空間計算量: O(1)
 *
 * @param arr - 検索対象のソート済み配列
 * @param target - 検索する要素
 * @returns 対象要素が見つかった場合はそのインデックス、見つからなかった場合は-1
 */
export function binarySearch<T>(arr: T[], target: T): number {
  if (arr.length === 0) return -1;

  let left = 0;
  let right = arr.length - 1;

  while (left <= right) {
    const mid = Math.floor((left + right) / 2);

    if (arr[mid] === target) {
      return mid;
    }

    if (arr[mid] < target) {
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }

  return -1;
}

/**
 * 二分探索を使って要素の挿入位置を返す
 * ソート順を維持するために要素が挿入されるべきインデックスを返す
 *
 * 時間計算量: O(log n)
 * 空間計算量: O(1)
 *
 * @param arr - ソート済み配列
 * @param target - 挿入位置を探す要素
 * @returns 要素が挿入されるべきインデックス
 */
export function binarySearchInsertionPoint<T>(arr: T[], target: T): number {
  if (arr.length === 0) return 0;

  let left = 0;
  let right = arr.length - 1;

  while (left <= right) {
    const mid = Math.floor((left + right) / 2);

    if (arr[mid] === target) {
      return mid;
    }

    if (arr[mid] < target) {
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }

  return left;
}

/**
 * カスタム比較関数を使用した二分探索
 * オブジェクトや複雑な型の配列に便利
 *
 * 時間計算量: O(log n)
 * 空間計算量: O(1)
 *
 * @param arr - ソート済み配列
 * @param target - 検索対象
 * @param compareFn - 要素を比較する関数。以下の値を返す必要がある:
 *                    a < b の場合は負の値、a === b の場合は0、a > b の場合は正の値
 * @returns 対象要素が見つかった場合はそのインデックス、見つからなかった場合は-1
 */
export function binarySearchWithComparator<T, U>(
  arr: T[],
  target: U,
  compareFn: (element: T, target: U) => number
): number {
  if (arr.length === 0) return -1;

  let left = 0;
  let right = arr.length - 1;

  while (left <= right) {
    const mid = Math.floor((left + right) / 2);
    const comparison = compareFn(arr[mid], target);

    if (comparison === 0) {
      return mid;
    }

    if (comparison < 0) {
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }

  return -1;
}
