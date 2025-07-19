 You will be tasked to.
  Generate a Motion Canvas 2D animation script (`.tsx` file) demonstrating a translation of a circle, a static LaTeX matrix, and a line. Follow the developer style guidelines exactly.

  🧑‍💻 Developer Style Guidelines (CRITICAL)

  1. **Dynamic Value Initialization & Dependencies**  
     - Use functions-as-values and `createComputed` for all numeric properties (`x`, `y`, `width`, `height`, `points`).  
     - Create reactive chains via `createRef` and computed properties so that updates cascade automatically.
     - Initialize the view with a black background color using `view.fill('#000000');`

  2. **Layout Paradigm (NO FLEXBOX)**  
     - ❌ Do not use `Layout`.  
     - ✅ Use `Rect` (center-anchored) as containers and `Node` for precise relative positioning.  
     - Construct a direct parent→child hierarchy only.

  3. **Relative Positioning**  
     - Compute positions with parent dimensions (`.width()`, `.height()`)—no hard-coded pixels.

  4. **External Utilities** 
    Avoid using the external utilities, as much as possible. eg.  
     ```ts
     import { logMethods, recurse_parent_with_width_height } from './utils';
     ```
  5. **Imports**
    (CRITICAL) 
    You must import all Motion-Canvas related code from '@motion-canvas/2d' or '@motion-canvas/core';
    eg.
    import {Circle, Layout, Rect, Node, makeScene2D, Txt, saturate, contrast} from '@motion-canvas/2d';
    import {
      all,
      createRef,
      easeInExpo,
      easeInOutExpo,
      waitFor,
      waitUntil,
      ThreadGenerator,
      chain,
      createSignal,
      slideTransition,
      Direction,
      easeOutCirc,
      createEaseInOutBack,
      range,
      InterpolationFunction
    } from '@motion-canvas/core';

  6. **OUTPUT**
    (CRITICAL) 
    You must ultimately output the written .tsx file to the frontend/src/scenes/example.tsx file. It already exists, but you must call the tool to edit it. 

I will provide you some exmaples of Motion Canvas Code, for you to see the syntax and these guidelines in practice. 

EXAMPLES :
