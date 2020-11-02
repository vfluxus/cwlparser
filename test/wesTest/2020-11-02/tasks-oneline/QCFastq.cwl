cwlVersion: v1.0

class: CommandLineTool 
id: qc_fastq 

requirements:
  - class: DockerRequirement
    dockerPull: vinbdi/biotools
  - class: ResourceRequirement
    ramMin: 4096
    coresMin: 4

inputs:
  input_fastqs:
    type: File[]
  

outputs:
  output: 
    type: File 
    outputBinding:
      glob: '*.zip'
      
arguments: 
  - position: 0
    shellQuote: false 
    valueFrom: >- 
      $(inputs.input_fastqs) -o .
baseCommand: [fastqc]  