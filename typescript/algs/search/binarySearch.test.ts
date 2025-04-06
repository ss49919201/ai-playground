/// <reference path="./jest.d.ts" />
import {
  binarySearch,
  binarySearchInsertionPoint,
  binarySearchWithComparator,
} from "./binarySearch";

describe("binarySearch", () => {
  test("ソート済み配列から要素を見つける", () => {
    expect(binarySearch([1, 2, 3, 4, 5], 3)).toBe(2);
    expect(binarySearch([1, 3, 5, 7, 9], 5)).toBe(2);
  });

  test("要素が見つからない場合は-1を返す", () => {
    expect(binarySearch([1, 2, 3, 4, 5], 6)).toBe(-1);
    expect(binarySearch([1, 3, 5, 7, 9], 4)).toBe(-1);
  });

  test("空の配列の場合は-1を返す", () => {
    expect(binarySearch([], 5)).toBe(-1);
  });

  test("配列の端にある要素を見つける", () => {
    expect(binarySearch([1, 2, 3, 4, 5], 1)).toBe(0);
    expect(binarySearch([1, 2, 3, 4, 5], 5)).toBe(4);
  });

  test("文字列の配列でも動作する", () => {
    expect(binarySearch(["a", "b", "c", "d", "e"], "c")).toBe(2);
    expect(binarySearch(["a", "b", "c", "d", "e"], "f")).toBe(-1);
  });
});

describe("binarySearchInsertionPoint", () => {
  test("要素の挿入位置を正しく見つける", () => {
    expect(binarySearchInsertionPoint([1, 3, 5, 7], 4)).toBe(2);
    expect(binarySearchInsertionPoint([1, 3, 5, 7], 0)).toBe(0);
    expect(binarySearchInsertionPoint([1, 3, 5, 7], 9)).toBe(4);
  });

  test("要素が存在する場合はその位置を返す", () => {
    expect(binarySearchInsertionPoint([1, 3, 5, 7], 3)).toBe(1);
    expect(binarySearchInsertionPoint([1, 3, 5, 7], 7)).toBe(3);
  });

  test("空の配列の場合は0を返す", () => {
    expect(binarySearchInsertionPoint([], 5)).toBe(0);
  });
});

describe("binarySearchWithComparator", () => {
  test("カスタム比較関数で要素を見つける", () => {
    const arr = [
      { id: 1, name: "Alice" },
      { id: 2, name: "Bob" },
      { id: 3, name: "Charlie" },
      { id: 4, name: "Dave" },
    ];

    // IDで検索
    const idComparator = (elem: { id: number }, target: number) =>
      elem.id - target;
    expect(binarySearchWithComparator(arr, 3, idComparator)).toBe(2);
    expect(binarySearchWithComparator(arr, 5, idComparator)).toBe(-1);

    // 名前で検索
    const nameComparator = (elem: { name: string }, target: string) => {
      if (elem.name < target) return -1;
      if (elem.name > target) return 1;
      return 0;
    };
    expect(binarySearchWithComparator(arr, "Bob", nameComparator)).toBe(1);
    expect(binarySearchWithComparator(arr, "Eve", nameComparator)).toBe(-1);
  });

  test("空の配列の場合は-1を返す", () => {
    expect(binarySearchWithComparator([], 5, () => 0)).toBe(-1);
  });
});
