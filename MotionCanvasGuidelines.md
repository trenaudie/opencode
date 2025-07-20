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
    import {Circle, Layout, Rect, Node, makeScene2D, Txt, SVG} from '@motion-canvas/2d';
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

    7. **Animation Flow**  
     - See the animations in the specs below. Avoid if possible to use tweens, as they are complex objects to code. 


I will provide you with a spec format for the Motion Canva code. A motion canvas scene should implement the features of the spec, but is not required to implement every feature. 
If you want to give instructions for generating Motion Canvas code, you MUST use the spec as a template for said instructions. 
I may also provide you some exmaples of Motion Canvas Code, for you to see the syntax and these guidelines in practice. 

## SPEC FORMAT : 
{
  "title": "string (scene short name/label)",
  "description": "string (concise, what is shown/animated, key visual/logic concept)",

  "sceneMetadata": {
    "background": "string|null",                                // Used if explicitly set, else omit/null
    "canvasDefaults": "object|null",                            // As present (empty/default if not specified)
    "viewport": "string|null",                                  // Only if present in scene
    "canvasResolution": "object|null",                          // Only if present
    "other": "object|null"                                      // Misc, if needed (from scene objects)
  },

  "imports": [
    {
      "type": "component|function|utility|reactive|utility/internal",
      "name": "string|array",
      "from": "string (import source module)"
    }
    // ...one object per import
  ],

  "objectInitialization": [
    // All nodes, layouts, blocks, visuals, in order (may be nested with 'children')
    {
      "name": "string (optional, for referencing/hierarchy clarity)",
      "type": "component/class (Circle, Layout, CodeBlock, Rect, Node, Path, Grid, Line, Txt, Latex, etc.)",
      "ref": "string|null",                             // e.g. "colA", "circleRefs[0]", or null
      "parent": "string|null",                          // For top-level, value is "view", or another name in this list
      "properties": {
        // Only those present across all examples, Partial, as used
        "x": "number|string|null",
        "y": "number|string|null",
        "position": "[number,number]|string|null",
        "size": "number|[number,number]|null",
        "width": "number|string|null",
        "height": "number|string|null",
        "fill": "string|null",
        "stroke": "string|null",
        "lineWidth": "number|null",
        "start": "number|null",
        "end": "number|null",
        "scale": "number|null",
        "radius": "number|null",
        "grow": "number|null",
        "gap": "number|null",
        "padding": "number|null",
        "margin": "number|null",
        "alignItems": "string|null",
        "justifyContent": "string|null",
        "direction": "string|null",
        "layout": "boolean|null",
        "fontWeight": "number|string|null",
        "fontSize": "number|null",
        "fontFamily": "string|null",
        "opacity": "number|null",
        "text": "string|array|null",
        "tex": "string|array|null",
        "code": "string|null",
        "data": "string|null",
        "arrowSize": "number|null",
        "endArrow": "boolean|null",
        "lineCap": "string|null",
        "points": "array|null",
        "others": "object|null"
      },
      "anchor": "string|null",
      "staticStyling": "object|string|null",
      "children": "array|null",    // (Repeat objectInitialization format, supports deep nesting)
      "count": "number|null"       // If this object is a group/array
    }
    // ...more as needed by scene
  ],

  "refs": {
    // (Optional map) Keyed by ref name or variable, value is node/component type or usage
    "refName": "type/description"
  },

  "relativePositioning": {
    "hierarchy": [
      // Explicit or compact textual/array tree, top-down
      // E.g. ["view", ["Layout", ["Rect(colA)", ...]]]
      // or, dictionaries for clarity:
      {
        "parent": "string",
        "children": ["string", ...]
      }
      // ...nested/composed as needed for clarity
    ],
    "layoutLogic": [
      "string (explain flex/fill/grow/column rules, per layout, as in your grid/flex scenes)",
      "string (note, e.g., 'colA and colB grow proportionally, rowA vertical in column, etc.')"
    ],
    "alignment": [
      "string (describe horizontal/vertical alignment, default behaviors, notes)"
    ],
    "placement": [
      // Provide additional location/mapping info, especially for animated or axis-projected objects
      "string (e.g., 'Tracker Y matches moving S&P point', 'Label at (x,y), ...')"
    ],
    "referenceSummary": [
      // If used, summary for named refs/variables, to correlate with code
      "string"
    ]
  },

  "animationFlow": [
    // Each animation step or phase as an object (order in array = time order in scene)
  {
    "step": "number|string",             // Sequence order or phase
    "type": "tween|spring|wait|custom",  // MC animation type

    // Targets
    "target": "string|array|null",       // Reference/name(s) of animated object(s)
    "property": "string|array|null",     // 'x', 'fill', ['x','y'], 'fillColor', etc

    // Value transitions
    "from": "number|string|array|object|null",   // Start value(s)
    "to": "number|string|array|object|null",     // End value(s)
    "startValue": "number|string|array|object|null", // Alias for compatibility
    "endValue": "number|string|array|object|null",

    // Duration / time
    "duration": "number|null",            // In seconds (support float/duration/named)
    "timingFunction": "string|null",      // Easing, 'easeInOutCubic' etc. Also can specify for springs.
    "easing": "string|null",              // Alias for compatibility
    "springType": "string|null",          // For spring animation ('PlopSpring', 'SmoothSpring', etc)

    // Tween/step logic
    "mapping": "string|null",             // e.g., 'map(start, end, easeFn(value))'
    "stepLogic": "string|object|null",    // Human summary, or actual mapped MC code for value updating
    "animationFn": "tween|spring|waitFor|edit|all|map|custom|null", // Explicit function if significant
    "call": "string|null",                // If 'callable', e.g., how the tween or spring is used

    // Control/concurrency
    "concurrency": "single|parallel|sequential|staggered|grouped|none|null", // Execution relation
    "runs": "description of sequence or concurrency",
    "order": "number|string|null",

    // Functional or advanced options
    "tweenCallback": "string|object|null",   // If a callback is used for custom value mapping in tween
    "easingCurve": "string|array|null",      // e.g. [custom array] if applicable
    "delays": "array|number|null",           // For staggered/concurrent tweens
    "labelsOrComments": "array|string|null", // Timeline comments, labels, debug highlights, as present

    // Waits and pauses
    "wait": "number|null",                   // Pause after/before, in seconds
    "waitType": "waitFor|yield|custom|null", // To clarify wait style if ambiguous

    // Additional
    "action": "string|null",                 // Human description ("fade out", "move left", etc)
    "details": "object|string|null",         // Any other MC-specific tweak, mapping, additional callback info.
    "notes": "string|null"                   // Freeform for implementer ("spring reverses direction", etc)
  }
  ],

  "signalsAndEffects": [
    // Only include if actually present in scene
    {
      "type": "Signal|Effect",
      "name": "string",
      "usage": "string (what/why/how used)"
    }
    // ...more as present
  ],

  "dataLogic": [
    // Only present in data-rich or chart/graph scenes
    {
      "syntheticData": {
        "strategy": "string|null",
        "avgReturn": "string|null",
        "volatility": "string|null",
        "weekendHandling": "string|null"
      },
      "length": "number|null",
      "mapping": {
        "X axis (time)": "string|null",
        "Y axis (value)": "string|null"
      }
    }
    // ...add logic notes as present
  ],

  "enhancements": [
    // Optional, only if present in scene: e.g. labels, tracker annotations, timeline comments
    {
      "type": "label|tracker|comment|annotation",
      "location": "string|null",
      "text": "string|array|null",
      "logic": "string|array|null"
    }
    // ...add more categories as needed
  ],

  "OptionalEnhancements": {
    "loops": "false|object|string|null",
    "conditionals": "false|object|string|null",
    "eventTriggers": "false|object|string|null",
    "interactivity": "false|object|string|null",
    "timelineCommentsOrLabels": "false|array|object|null",
    "timelineAnnotations": "array|object|null",
    "consoleLogging": "array|object|null",
    "sceneMetadata": "object|null"
  },

  "notes": [
    // Freeform, add any additional clarification needed for recreation
    "string"
  ]
}

## Example of a perfect scene
Here is a scene that implements these guidelines and compiles successfully. 
It is a good representation of how to use SVGs.
The SVGs must be imported as a raw string from the public dir, using the '?raw' flag
Always use the following format import logo from '/public/<svgname>.svg?raw';
eg. import logo from '/public/logo.svg?raw';
CRITICAL: You must select the Path children of the SVG and fill them with a color. Use a FOR loop for this. See below.
import {
  Rect,
  Node,
  makeScene2D, SVG,Path
} from '@motion-canvas/2d';
import {
  createSignal,
  createRef,
  all,
  waitFor,
  linear,
  Vector2
} from '@motion-canvas/core';
import logo from '/public/logo.svg?raw';
console.log(`found logo ${logo}`)
const numRows = 5;
const baseWidth = 560;
const widthStep = 80;
const boxHeight = 48;
const gap = 8;

export default makeScene2D(function* (view) {
  // Set background color
  view.fill('#000000');

  // Refs for rectangles and signals for opacity
  const rectRefs = Array.from({ length: numRows }, () => createRef<Rect>());
  const rectSignals = Array.from({ length: numRows }, () => createSignal(0));
  const logoRef = createRef<SVG>();
  // Main container Node to stack rectangles
  const container = createRef<Node>();
  const gap_between_last_rectangle_and_svg_y = 250;
  view.add(
    <>
  <Node ref={container}>
      {Array.from({ length: numRows }, (_, i) => {
        const width = baseWidth - i * widthStep;
        const y = () => (i * (boxHeight + gap)) - (numRows * (boxHeight + gap) / 2) + (boxHeight / 2);

        return (
          <Rect
            ref={rectRefs[i]}
            width={width}
            height={boxHeight}
            fill={'#6cf1c2'}
            y={y}
            opacity={rectSignals[i]} // opacity controlled by signal
          />
        );
      })}
    </Node>
      <SVG ref={logoRef} svg={logo.replace("@color", "#f2ff48")} size = {300}  position= {() => {
        let svg_world_to_parent = logoRef().worldToParent();
      console.log(`the last child of the container is ${container().children().at(-1).absolutePosition()}`);
            console.log(`which, converting to the Parent of the SVG means: ${svg_world_to_parent.transformPoint(container().children().at(-1).absolutePosition())}`);
            let gap_between_last_rectangle_and_svg_vector =  new Vector2(0, gap_between_last_rectangle_and_svg_y);
            console.log(`gap is ${gap_between_last_rectangle_and_svg_vector}`)
    return container().children().at(-1).absolutePosition().transformAsPoint(svg_world_to_parent).add(gap_between_last_rectangle_and_svg_vector)
    }}/>
  </>
    
  );
  // Animate SVG elements
  const svgElements = [logoRef()]
  const animations = [];
  for (const svgElement of svgElements) {
    animations.push(svgElement.scale(1, .5));
    animations.push(svgElement.opacity(1, .5));
    for (const child of svgElement.children()[0].children()) {
      if (child instanceof Path) {
        yield* child.fill('white',1);
      }
    }
  }
  // Fade in rectangles one by one
  for (let i = 0; i < numRows; i++) {
    yield* all(
      rectSignals[i](1, 0.5, linear), // Fade to opaque
      waitFor(0.2) // Wait before next rectangle
    );
  }
});


It is also a good representation of how you can position points one to another, you must always give the x={} and y={} position relative to the PARENT element. So, for example, if you want object A to have a position that depends on object B that is not the parent, one option is:
1. to fetch the ABSOLUTE position of object B. eg. objectBAbsPos = objectBRef().absolutePosition()
2. convert it to the local coordinates of the PARENT of object A. eg. objectBAbsPos_inParentcoords = objectBAbsPos.transformAsPoint(objectARef().worldToParent())
3. add the offset that you want to create between the two objects. eg. objectAPos = objectBAbsPos_inParentcoords.add(new Vector2([0,200]))
