export interface User {}

export interface DB {
  getUsers(input: { ids: number[]; limit: number }): User[];
}

export const newDB = (): DB => ({
  getUsers(input: { ids: number[]; limit: number }): User[] {
    return [];
  },
});

export const newSearchUserUsecase = (db: DB) => ({
  exec(input: { ids: number[]; limit: number }) {
    return db.getUsers(input);
  },
});

interface DIContainer {
  get<T>(key: string): T;
  register<T>(key: string, factory: () => T): void;
}

export const createContainer = (): DIContainer => {
  const services = new Map<string, () => unknown>();

  return {
    register<T>(key: string, factory: () => T): void {
      services.set(key, factory);
    },
    get<T>(key: string): T {
      const factory = services.get(key);
      if (!factory) {
        throw new Error(`Service not found: ${key}`);
      }
      return factory() as T;
    },
  };
};

const container = createContainer();
container.register('db', () => newDB());
container.register('searchUserUsecase', () => newSearchUserUsecase(container.get<DB>('db')));

export const SearchUserHandler = (input: { ids: number[]; limit: number }) => {
  const usecase = container.get<ReturnType<typeof newSearchUserUsecase>>('searchUserUsecase');
  return usecase.exec(input);
};
