import { Exercise, Set, TrainingRecord } from "../types";

let trainingRecords: TrainingRecord[] = [];

export const generateId = (): string => {
  return Math.random().toString(36).substring(2, 9);
};

export const getCurrentISOString = (): string => {
  return new Date().toISOString();
};

export const addTrainingRecord = (
  record: Omit<TrainingRecord, "id" | "createdAt" | "updatedAt">
): TrainingRecord => {
  const timestamp = getCurrentISOString();
  const newRecord: TrainingRecord = {
    ...record,
    id: generateId(),
    createdAt: timestamp,
    updatedAt: timestamp,
  };

  trainingRecords = [...trainingRecords, newRecord];
  return newRecord;
};

export const getTrainingRecord = (
  id: string
): TrainingRecord | undefined => {
  return trainingRecords.find((record) => record.id === id);
};

export const getAllTrainingRecords = (): TrainingRecord[] => {
  return trainingRecords;
};

export const updateTrainingRecord = (
  id: string,
  updates: Partial<Omit<TrainingRecord, "id" | "createdAt" | "updatedAt">>
): TrainingRecord | undefined => {
  let updatedRecord: TrainingRecord | undefined;
  
  trainingRecords = trainingRecords.map((record) => {
    if (record.id === id) {
      updatedRecord = {
        ...record,
        ...updates,
        updatedAt: getCurrentISOString(),
      };
      return updatedRecord;
    }
    return record;
  });

  return updatedRecord;
};

export const deleteTrainingRecord = (id: string): boolean => {
  const initialLength = trainingRecords.length;
  trainingRecords = trainingRecords.filter((record) => record.id !== id);
  return trainingRecords.length !== initialLength;
};

export const addExerciseToRecord = (
  recordId: string,
  exercise: Omit<Exercise, "id">
): Exercise | undefined => {
  let newExercise: Exercise | undefined;
  
  trainingRecords = trainingRecords.map((record) => {
    if (record.id === recordId) {
      newExercise = {
        ...exercise,
        id: generateId(),
      };
      
      return {
        ...record,
        exercises: [...record.exercises, newExercise],
        updatedAt: getCurrentISOString(),
      };
    }
    return record;
  });

  return newExercise;
};

export const addSetToExercise = (
  recordId: string,
  exerciseId: string,
  set: Omit<Set, "id">
): Set | undefined => {
  let newSet: Set | undefined;
  
  trainingRecords = trainingRecords.map((record) => {
    if (record.id === recordId) {
      const updatedExercises = record.exercises.map((exercise) => {
        if (exercise.id === exerciseId) {
          newSet = {
            ...set,
            id: generateId(),
          };
          
          return {
            ...exercise,
            sets: [...exercise.sets, newSet],
          };
        }
        return exercise;
      });

      return {
        ...record,
        exercises: updatedExercises,
        updatedAt: getCurrentISOString(),
      };
    }
    return record;
  });

  return newSet;
};
