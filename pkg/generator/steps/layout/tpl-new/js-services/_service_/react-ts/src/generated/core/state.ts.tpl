// vim: set ft=typescript:
import { createSlice } from '@reduxjs/toolkit';
import MifyContext from './context';

export interface MifyContextState {
  value: MifyContext;
}

const initialState: MifyContextState = {
  value: new MifyContext(),
};

export const mifyContextSlice = createSlice({
  name: 'mify_context',
  initialState,
  reducers: {},
  extraReducers: () => {},
});

export default mifyContextSlice.reducer;
