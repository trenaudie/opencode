You're given a TypeScript Motion Canvas file. Your task is to analyze and describe its structure and behavior as a detailed technical specification. You may invent your own spec structure (in JSON though), but it must be structured and exhaustive. The description must cover the following five key aspects:
. Imports
List all external and internal modules, components, or assets imported at the top of the file.

Include import types (e.g., component, utility, asset, type) and the source path.

2. Object Initialization
Describe every graphical or animated object instantiated (e.g., Rect, Txt, Img, Node, Circle, SVG, etc.).

For each object, specify:

Type/class of the object

Initial properties (e.g., position, size, color, opacity, rotation)

Anchor or origin point

Any static styling (e.g., border radius, shadows, fills)

Parenting structure (i.e., if it is a child of another object or part of a group)

3. Relative Positioning
Explain the spatial layout and containment hierarchy:

Which objects are nested inside others (parent-child relationships)?

What are the positions of objects relative to each other (e.g., center-aligned, above, inside, offset by a margin)?

Describe any layout logic or utilities used for placement.

4. Animation Flow
Document the full sequence of animations performed in the scene:

For each animation or transition, include:

The target object

The property being animated (e.g., opacity, position, scale, rotation)

The animation function (e.g., tween, waitFor, chain, all)

Start and end values, where applicable

Duration and easing function, if defined

Clarify the order and concurrency of animations (i.e., what runs in parallel vs sequentially).

5. Optional Enhancements
If present, also include:

Any looping logic, conditionals, or event triggers

Timeline comments or labels

Interactive components or user-triggered transitions

Scene-level metadata or configuration (e.g., canvas resolution, background settings)


💡 Goal: Your spec should be precise enough for someone to recreate the scene without seeing the original code.

You can add other elements to the spec if you wish. 
