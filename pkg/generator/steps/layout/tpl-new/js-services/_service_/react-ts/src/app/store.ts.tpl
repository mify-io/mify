// vim: set ft=typescript:
import { configureStore, ThunkAction, Action } from '@reduxjs/toolkit';
import mifyContextReducer from '../generated/core/state';

export const store = configureStore({
  reducer: {
    mifyState: mifyContextReducer,
  },
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
