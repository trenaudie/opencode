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
      "step": "number|string",
      "type": "sequential|parallel|staggered|wait|custom",
      "target": "string|array|null",                     // Object(s) acted upon (often by 'ref')
      "property": "string|array|null",                   // Animated property(ies)
      "from": "value|object|null",
      "to": "value|object|null",
      "duration": "number|null",                         // In seconds if present
      "easing": "string|null",
      "animationFn": "tween|all|waitFor|edit|insert|remove|selection|custom|null",
      "action": "string|null",                           // Plaintext summary of "what happens"
      "details": "object|string|null",                   // Extra details: pausing, method calls, value calcs.
      "concurrency": "sequential|parallel|staggered|null"
    }
    // ...sequence continues per scene animation script
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