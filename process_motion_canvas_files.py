#!/usr/bin/env python3

import os
import json
from pathlib import Path
from openai import OpenAI
import itertools
def process_tsx_files():
    """
    Process .tsx files in frontend/src/scenes2 and frontend/src/scenes3 directories.
    Skip files that already have corresponding _spec.json files.
    """
    
    # Initialize OpenAI client
    client = OpenAI()
    
    # Define directories to process
    directories = [
        "frontend/src/scenes2",
        "frontend/src/scenes3"
    ]
    
    # Files to skip (already processed)
    skip_files = {
        "example2.tsx",
        "example3.tsx", 
        "example4.tsx"
    }
    
    for directory in directories:
        if not os.path.exists(directory):
            print(f"Directory {directory} does not exist, skipping...")
            continue
            
        print(f"Processing directory: {directory}")
        
        # Get all .tsx files in the directory
        tsx_files = [f for f in os.listdir(directory) 
                    if f.endswith('.tsx') and f not in skip_files]
        
        for tsx_file in itertools.islice(tsx_files, None):  # Limit to first 10 files for processing
            file_path = os.path.join(directory, tsx_file)
            
            # Generate output filename
            base_name = tsx_file.replace('.tsx', '')
            output_file = os.path.join(directory, f"{base_name}_spec.json")
            
            # # Skip if spec file already exists
            # if os.path.exists(output_file):
            #     print(f"Spec file already exists for {tsx_file}, skipping...")
            #     continue
            
            try:
                print(f"Processing {tsx_file}...")
                
                # Read the .tsx file content
                with open(file_path, 'r', encoding='utf-8') as f:
                    tsx_content = f.read()
                
                # Call OpenAI API with the specified prompt
                response = client.chat.completions.create(
                    model="gpt-4.1",
                    messages=[
                        {"role": "system", "content": Path("trenaudie_prompts/prompt_summarize_motion_canvas_scene_2.md").read_text(encoding='utf-8')},
                        {"role": "user", "content": f"Please summarize the following Motion Canvas scene:\n\n{tsx_content}"}
                    ]
)
                # Save the response to the spec file
                spec_data = {
                    "source_file": tsx_file,
                    "source_path": file_path,
                    "openai_response": response.choices[0].message.content.strip(),
                }
                
                with open(output_file, 'w', encoding='utf-8') as f:
                    json.dump(spec_data, f, indent=2, ensure_ascii=False)
                
                print(f"✅ Successfully processed {tsx_file} -> {base_name}_spec.json")
                
            except Exception as e:
                print(f"❌ Error processing {tsx_file}: {str(e)}")
                continue

def main():
    """Main function to run the script."""
    print("🚀 Starting Motion Canvas file processing...")
    print("=" * 50)
    
    # Check if we're in the right directory
    if not os.path.exists("frontend"):
        print("❌ Error: 'frontend' directory not found.")
        print("Please run this script from the project root directory.")
        return 1
    
    try:
        process_tsx_files()
        print("=" * 50)
        print("✅ Processing complete!")
        return 0
        
    except Exception as e:
        print(f"❌ Fatal error: {str(e)}")
        return 1

if __name__ == "__main__":
    exit(main())